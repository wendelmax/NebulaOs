package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/wendelmax/nebulaos/src/api/domain"
	"github.com/wendelmax/nebulaos/src/api/internal/api"
	"github.com/wendelmax/nebulaos/src/api/internal/api/middleware"
	"github.com/wendelmax/nebulaos/src/api/internal/infrastructure"
	"github.com/wendelmax/nebulaos/src/api/internal/services"
	"github.com/wendelmax/nebulaos/src/api/internal/usecase"
	"github.com/wendelmax/nebulaos/src/providers"
	"github.com/wendelmax/nebulaos/src/providers/aws"
	"github.com/wendelmax/nebulaos/src/providers/baremetal"
	"github.com/wendelmax/nebulaos/src/providers/keycloak"
	"github.com/wendelmax/nebulaos/src/providers/kubernetes"
	"github.com/wendelmax/nebulaos/src/providers/mock"
	"github.com/wendelmax/nebulaos/src/providers/openstack"
	"github.com/wendelmax/nebulaos/src/providers/proxmox"
	"github.com/wendelmax/nebulaos/src/providers/storage"
	"github.com/wendelmax/nebulaos/src/providers/traefik"
	"github.com/wendelmax/nebulaos/src/providers/vault"
	_ "github.com/lib/pq" // NEW
	"github.com/nats-io/nats.go"
)

func main() {
	fmt.Println("Starting NebulaOS Cloud API (Phase 11: Production Hardened)...")

	// Infrastructure Self-Healing
	infrastructureSelfHeal()

	// Configuration
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	kcURL := getEnv("KC_URL", "http://localhost:8080")
	vaultURL := getEnv("VAULT_URL", "http://localhost:8200")
	vaultToken := getEnv("VAULT_TOKEN", "root-token")
	dbURL := os.Getenv("DATABASE_URL")

	// Infrastructure - NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to NATS at %s: %v. Audit logging might be degraded.", natsURL, err)
	} else {
		defer nc.Close()
		fmt.Printf("Connected to NATS at %s for Audit Logging\n", natsURL)
	}

	// Dependency Injection - Providers
	providerFactory := providers.NewProviderFactory()
	storageProvider := storage.NewMockStorageProvider()

	mockProvider := mock.NewMockProvider()
	providerFactory.Register("mock", mockProvider)

	proxmoxURL := getEnv("PROXMOX_URL", "https://pve.nebula.local/api2/json")
	proxmoxToken := getEnv("PROXMOX_TOKEN", "token-uuid")
	proxmoxProvider := proxmox.NewProxmoxProvider(proxmoxURL, proxmoxToken)
	providerFactory.Register("proxmox", proxmoxProvider)

	kubeconfig := getEnv("KUBECONFIG_PATH", "kubeconfig-data")
	k8sProvider := kubernetes.NewKubernetesProvider(kubeconfig)
	providerFactory.Register("kubernetes", k8sProvider)

	osEndpoint := getEnv("OS_IDENTITY_ENDPOINT", "http://openstack:5000/v3")
	osUser := getEnv("OS_USERNAME", "admin")
	osPass := getEnv("OS_PASSWORD", "password")
	osProvider, _ := openstack.NewOpenStackProvider(gophercloud.AuthOptions{
		IdentityEndpoint: osEndpoint,
		Username:         osUser,
		Password:         osPass,
	})
	if osProvider != nil {
		providerFactory.Register("openstack", osProvider)
	}

	bmSSHHost := getEnv("BM_SSH_HOST", "localhost")
	bmSSHUser := getEnv("BM_SSH_USER", "root")
	bmSSHKey := getEnv("BM_SSH_KEY_PATH", "")
	bmSSHPass := getEnv("BM_SSH_PASSWORD", "")
	bmIPMIUser := getEnv("BM_IPMI_USER", "admin")
	bmIPMIPass := getEnv("BM_IPMI_PASSWORD", "password")
	bmProvider := baremetal.NewBareMetalProvider(bmSSHHost, bmSSHUser, bmSSHKey, bmSSHPass, bmIPMIUser, bmIPMIPass)
	providerFactory.Register("baremetal", bmProvider)

	awsRegion := getEnv("AWS_REGION", "us-east-1")
	awsEndpoint := getEnv("AWS_ENDPOINT", "http://localhost:4566")
	awsProvider, _ := aws.NewAWSProvider(context.Background(), awsRegion, awsEndpoint)
	if awsProvider != nil {
		providerFactory.Register("aws", awsProvider)
	}

	traefikConfigDir := getEnv("TRAEFIK_CONFIG_DIR", "./configs/traefik")
	traefikProvider := traefik.NewTraefikProvider(traefikConfigDir)
	kcClientID := getEnv("KC_CLIENT_ID", "nebula-api")
	keycloakProvider := keycloak.NewKeycloakProvider(kcURL, kcClientID)
	vaultProvider := vault.NewVaultProvider(vaultURL, vaultToken)


	// Dependency Injection - Repositories
	var tenantRepo domain.TenantRepository
	var projectRepo domain.ProjectRepository
	var resourceRepo domain.ResourceRepository
	var quotaRepo domain.QuotaRepository
	var volumeRepo domain.VolumeRepository
	var bucketRepo domain.BucketRepository
	var policyRepo domain.SovereigntyPolicyRepository
	var securityGroupRepo domain.SecurityGroupRepository
	var tfStateRepo domain.TerraformStateRepository
	var blueprintRepo domain.BlueprintRepository
	var gslbRepo domain.GSLBRepository
	var regionRepo domain.RegionRepository
	var zoneRepo domain.AvailabilityZoneRepository
	var providerRepo domain.ProviderRepository
	var userRepo domain.UserRepository
	var orgRepo domain.OrganizationRepository
	var deptRepo domain.DepartmentRepository
	var bmRepo domain.BareMetalRepository

	var db *sql.DB
	if dbURL != "" {
		fmt.Println("Initializing PostgreSQL persistence layer...")
		var err error
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		
		// Run database bootstrap (Self-Healing Schema)
		if err := infrastructure.EnsureSchema(db); err != nil {
			log.Printf("Warning: Failed to ensure schema: %v. The application might be using an outdated database.", err)
		}

		tenantRepo = infrastructure.NewPostgresTenantRepository(db)
		projectRepo = infrastructure.NewPostgresProjectRepository(db)
		resourceRepo = infrastructure.NewPostgresResourceRepository(db)
		quotaRepo = infrastructure.NewPostgresQuotaRepository(db)
		volumeRepo = infrastructure.NewPostgresVolumeRepository(db)
		bucketRepo = infrastructure.NewPostgresBucketRepository(db)
		policyRepo = infrastructure.NewPostgresPolicyRepository(db)
		securityGroupRepo = infrastructure.NewPostgresSecurityGroupRepository(db)
		tfStateRepo = infrastructure.NewPostgresTerraformStateRepository(db)
		blueprintRepo = infrastructure.NewPostgresBlueprintRepository(db)
		gslbRepo = infrastructure.NewPostgresGSLBRepository(db)
		providerRepo = infrastructure.NewInMemoryProviderRepository() // Fallback or implement Postgres
		userRepo = infrastructure.NewPostgresUserRepository(db)
		orgRepo = infrastructure.NewPostgresOrganizationRepository(db)
		deptRepo = infrastructure.NewPostgresDepartmentRepository(db)
		bmRepo = infrastructure.NewPostgresBareMetalRepository(db)
	} else {
		fmt.Println("Initializing In-Memory repositories (Development Mode)...")
		tenantRepo = infrastructure.NewInMemoryTenantRepository()
		projectRepo = infrastructure.NewInMemoryProjectRepository()
		resourceRepo = infrastructure.NewInMemoryResourceRepository()
		quotaRepo = infrastructure.NewInMemoryQuotaRepository()
		volumeRepo = infrastructure.NewInMemoryVolumeRepository()
		bucketRepo = infrastructure.NewInMemoryBucketRepository()
		providerRepo = infrastructure.NewInMemoryProviderRepository()
		userRepo = infrastructure.NewInMemoryUserRepository()
		orgRepo = infrastructure.NewInMemoryOrganizationRepository()
		deptRepo = infrastructure.NewInMemoryDepartmentRepository()
		bmRepo = infrastructure.NewInMemoryBareMetalRepository()

		// Initial data for in-memory setup
		seedCtx := context.Background()
		if err := regionRepo.Create(seedCtx, &domain.Region{ID: getEnv("SEED_REGION_US_ID", "reg-us-east"), Name: getEnv("SEED_REGION_US_NAME", "US East (Proxmox Cluster A)"), Location: getEnv("SEED_REGION_US_LOCATION", "Virginia"), IsDefault: true}); err != nil {
			log.Printf("Warning: failed to seed US region: %v", err)
		}
		if err := regionRepo.Create(seedCtx, &domain.Region{ID: getEnv("SEED_REGION_EU_ID", "reg-eu-west"), Name: getEnv("SEED_REGION_EU_NAME", "EU West (OpenStack Lab)"), Location: getEnv("SEED_REGION_EU_LOCATION", "London")}); err != nil {
			log.Printf("Warning: failed to seed EU region: %v", err)
		}
		if err := zoneRepo.Create(seedCtx, &domain.AvailabilityZone{ID: getEnv("SEED_ZONE_ID", "zone-a"), RegionID: getEnv("SEED_REGION_US_ID", "reg-us-east"), Name: getEnv("SEED_ZONE_NAME", "Zone A"), State: "available"}); err != nil {
			log.Printf("Warning: failed to seed zone: %v", err)
		}

		// Seed Providers
		if err := providerRepo.Create(seedCtx, &domain.Provider{ID: getEnv("SEED_PROVIDER_PX_ID", "prov-px-1"), Name: "Proxmox Production", Type: domain.ProxmoxProvider, Endpoint: proxmoxURL, Status: "active", CreatedAt: time.Now()}); err != nil {
			log.Printf("Warning: failed to seed proxmox provider: %v", err)
		}
		if err := providerRepo.Create(seedCtx, &domain.Provider{ID: getEnv("SEED_PROVIDER_OS_ID", "prov-os-1"), Name: "OpenStack Lab", Type: domain.OpenStackProvider, Endpoint: osEndpoint, Status: "active", CreatedAt: time.Now()}); err != nil {
			log.Printf("Warning: failed to seed openstack provider: %v", err)
		}
	}

	// Seed internal admin user, defaults and test data
	seedEnterpriseDefaults(tenantRepo, projectRepo, orgRepo, deptRepo, userRepo, resourceRepo, volumeRepo, bucketRepo)

	fmt.Printf("Repositories initialized (Persistence: %v)\n", dbURL != "")


	// Choose Identity Manager (Keycloak or Internal)
	var identityManager domain.IdentityManager
	if os.Getenv("USE_KEYCLOAK") == "true" {
		identityManager = keycloakProvider
	} else {
		// Default to Internal Identity Manager
		jwtSecret := getEnv("JWT_SECRET", "nebula-secret-key-2026")
		identityManager = services.NewInternalIdentityManager(userRepo, jwtSecret)
	}

	authMiddleware := middleware.NewAuthMiddleware(identityManager)
	auditMiddleware := middleware.NewAuditMiddleware(nc)
	metricsMiddleware := middleware.NewMetricsMiddleware()

	// Services
	policyService := services.NewSovereignPolicyService(policyRepo)
	billingMgr := infrastructure.NewSovereignBillingManager(resourceRepo, volumeRepo, bucketRepo, tenantRepo)
	gslbManager := services.NewGSLBManager(gslbRepo)
	aiAdvisor := services.NewAIResourceAdvisor(resourceRepo)
	bmManager := services.NewBareMetalManager(bmRepo)
	networkMgr := services.NewNetworkManager(providerFactory)
	orchestrator := services.NewMetaOrchestrator(providerFactory, networkMgr)

	// Dependency Injection - Use Cases
	createTenantUC := usecase.NewCreateTenantUseCase(tenantRepo)
	createProjectUC := usecase.NewCreateProjectUseCase(projectRepo)
	getTenantUC := usecase.NewGetTenantUseCase(tenantRepo)
	listTenantsUC := usecase.NewListTenantsUseCase(tenantRepo)

	getProjectUC := usecase.NewGetProjectUseCase(projectRepo)
	listProjectsUC := usecase.NewListProjectsByTenantUseCase(projectRepo)

	// Networking Use Cases
	createSGUC := usecase.NewCreateSecurityGroupUseCase(securityGroupRepo)
	listSGsUC := usecase.NewListSecurityGroupsUseCase(securityGroupRepo)
	addFirewallRuleUC := usecase.NewAddFirewallRuleUseCase(securityGroupRepo)

	// Automation Use Cases
	saveTFStateUC := usecase.NewSaveTerraformStateUseCase(tfStateRepo)
	getTFStateUC := usecase.NewGetTerraformStateUseCase(tfStateRepo)
	listBlueprintsUC := usecase.NewListBlueprintsUseCase(blueprintRepo)
	createBlueprintUC := usecase.NewCreateBlueprintUseCase(blueprintRepo)

	createResourceUC := usecase.NewCreateResourceUseCase(resourceRepo, projectRepo, quotaRepo, policyService, providerFactory)
	deployBlueprintUC := usecase.NewDeployBlueprintUseCase(blueprintRepo, resourceRepo, providerFactory)
	autoProvisionUC := usecase.NewAutomatedProvisioningUseCase(deployBlueprintUC, blueprintRepo)
	getResourceUC := usecase.NewGetResourceUseCase(resourceRepo)
	listResourcesUC := usecase.NewListResourcesByProjectUseCase(resourceRepo)

	createVolumeUC := usecase.NewCreateVolumeUseCase(volumeRepo, storageProvider)
	createBucketUC := usecase.NewCreateBucketUseCase(bucketRepo, storageProvider)
	listVolumesUC := usecase.NewListVolumesUseCase(volumeRepo)
	listBucketsUC := usecase.NewListBucketsUseCase(bucketRepo)

	requestCertUC := usecase.NewRequestCertificateUseCase(traefikProvider)
	storeSecretUC := usecase.NewStoreSecretUseCase(vaultProvider)
	complianceUC := usecase.NewGetComplianceReportUseCase(resourceRepo, projectRepo, quotaRepo)

	// Dependency Injection - Handlers
	tenantHandler := api.NewTenantHandler(createTenantUC, getTenantUC, listTenantsUC)
	projectHandler := api.NewProjectHandler(createProjectUC, getProjectUC, listProjectsUC)
	resourceHandler := api.NewResourceHandler(createResourceUC, getResourceUC, listResourcesUC)
	storageHandler := api.NewStorageHandler(createVolumeUC, createBucketUC, listVolumesUC, listBucketsUC)
	bmHandler := api.NewBareMetalHandler(bmManager)
	hierarchyHandler := api.NewHierarchyHandler(orgRepo, deptRepo)
	complianceHandler := api.NewComplianceHandler(complianceUC)
	billingHandler := api.NewBillingHandler(billingMgr)
	policyHandler := api.NewPolicyHandler(policyService)

	// Routes
	mux := http.NewServeMux()

	// Auth Routes
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&input)
		token, err := identityManager.Authenticate(r.Context(), input.Username, input.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "nebula_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})

	mux.HandleFunc("/auth/me", func(w http.ResponseWriter, r *http.Request) {
		user, err := authMiddleware.AuthenticateRequest(r)
		if err != nil {
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "nebula_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		})
		json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
	})

	mux.Handle("/auth/change-password", authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		user := r.Context().Value("user").(*domain.User)
		var input struct {
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
			Email       string `json:"email"`
		}
		json.NewDecoder(r.Body).Decode(&input)
		
		// Internal implementation specific call
		if internalIM, ok := identityManager.(*services.InternalIdentityManager); ok {
			err := internalIM.ChangePassword(r.Context(), user.ID, input.OldPassword, input.NewPassword, input.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
		} else {
			http.Error(w, "Change password not supported for this identity provider", http.StatusNotImplemented)
		}
	})))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		status := map[string]interface{}{
			"api":     "active",
			"version": getEnv("API_VERSION", "v14.2-prod-hardened"),
			"time":    time.Now().Format(time.RFC3339),
		}
		infraHealthy := true

		// Check Postgres
		if db != nil {
			if err := db.PingContext(ctx); err != nil {
				status["database"] = "unhealthy: " + err.Error()
				infraHealthy = false
			} else {
				status["database"] = "healthy"
			}
		} else {
			status["database"] = "in-memory"
		}

		// Check NATS
		if nc != nil && nc.IsConnected() {
			status["nats"] = "healthy"
		} else {
			status["nats"] = "disconnected"
			// Audit might be degraded, but API can still serve if NATS is optional
		}

		// Check Keycloak
		if err := keycloakProvider.Ping(ctx); err != nil {
			status["keycloak"] = "unhealthy: " + err.Error()
			infraHealthy = false
		} else {
			status["keycloak"] = "healthy"
		}

		// Check Vault
		if err := vaultProvider.Ping(ctx); err != nil {
			status["vault"] = "unhealthy: " + err.Error()
			infraHealthy = false
		} else {
			status["vault"] = "healthy"
		}

		status["healthy"] = infraHealthy

		w.Header().Set("Content-Type", "application/json")
		if !infraHealthy {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(status)
	})

	// Network & Security Routes
	mux.Handle("/network/certificate", authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainName := r.URL.Query().Get("domain")
		if err := requestCertUC.Execute(r.Context(), domainName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Certificate requested"})
	})))

	// Secret Management
	mux.Handle("/secrets", authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var input struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		json.NewDecoder(r.Body).Decode(&input)
		if err := storeSecretUC.Execute(r.Context(), input.Key, input.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Secret stored"})
	})))

	// Networking & Security Groups
	mux.Handle("/security-groups", auditMiddleware.Audit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var input usecase.CreateSecurityGroupInput
			json.NewDecoder(r.Body).Decode(&input)
			if err := createSGUC.Execute(r.Context(), input); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Security group created"})
		} else {
			projectID := r.URL.Query().Get("project_id")
			sgs, err := listSGsUC.Execute(r.Context(), projectID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(sgs)
		}
	})))

	mux.Handle("/security-groups/rules", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var input usecase.AddFirewallRuleInput
			json.NewDecoder(r.Body).Decode(&input)
			if err := addFirewallRuleUC.Execute(r.Context(), input); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "Firewall rule added"})
		}
	}))))

	// Core API Routes
	mux.Handle("/tenants", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			tenantHandler.CreateTenant(w, r)
		} else {
			tenantHandler.ListTenants(w, r)
		}
	}))))

	mux.Handle("/projects", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			projectHandler.CreateProject(w, r)
		} else {
			projectHandler.ListProjects(w, r)
		}
	}))))

	mux.Handle("/resources", auditMiddleware.Audit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			resourceHandler.CreateResource(w, r)
		} else {
			resourceHandler.ListResources(w, r)
		}
	})))

	mux.Handle("/storage/volumes", auditMiddleware.Audit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			storageHandler.CreateVolume(w, r)
		} else {
			storageHandler.ListVolumes(w, r)
		}
	})))

	mux.Handle("/storage/buckets", auditMiddleware.Audit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			storageHandler.CreateBucket(w, r)
		} else {
			storageHandler.ListBuckets(w, r)
		}
	})))

	// Billing & Governance
	mux.Handle("/billing/report", auditMiddleware.Audit(http.HandlerFunc(billingHandler.GetReport)))
	mux.Handle("/compliance/report", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(complianceHandler.GetReport))))
	mux.Handle("/governance/policy", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			policyHandler.UpdatePolicy(w, r)
		} else {
			policyHandler.GetPolicy(w, r)
		}
	}))))

	// Automation: Terraform State
	mux.Handle("/automation/tf-state", authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID := r.URL.Query().Get("project_id")
		if r.Method == http.MethodPost {
			var input usecase.SaveTerraformStateInput
			json.NewDecoder(r.Body).Decode(&input)
			if err := saveTFStateUC.Execute(r.Context(), input); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "State saved"})
		} else {
			state, err := getTFStateUC.Execute(r.Context(), projectID)
			if err != nil {
				http.Error(w, "State not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(state)
		}
	})))

	// Marketplace: Blueprints
	mux.Handle("/marketplace/blueprints", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var b domain.Blueprint
			json.NewDecoder(r.Body).Decode(&b)
			if err := createBlueprintUC.Execute(r.Context(), &b); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Blueprint created"})
		} else {
			blueprints, _ := listBlueprintsUC.Execute(r.Context())
			json.NewEncoder(w).Encode(blueprints)
		}
	}))

	mux.Handle("/marketplace/deploy", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var input usecase.DeployBlueprintInput
		json.NewDecoder(r.Body).Decode(&input)
		if err := deployBlueprintUC.Execute(r.Context(), input); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Deployment initiated"})
	}))

	// Phase 14: Global Orchestration & AI
	mux.Handle("/network/gslb", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var ep domain.GlobalEndpoint
			json.NewDecoder(r.Body).Decode(&ep)
			if ep.ID == "" {
				ep.ID = domain.NewID()
			}
			if err := gslbManager.CreateEndpoint(r.Context(), &ep); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(ep)
		} else {
			eps, _ := gslbManager.ListEndpoints(r.Context())
			json.NewEncoder(w).Encode(eps)
		}
	}))

	mux.Handle("/intelligence/advisor", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID := r.URL.Query().Get("project_id")
		insights, err := aiAdvisor.AnalyzeUsage(r.Context(), projectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(insights)
	}))

	mux.Handle("/intelligence/stats", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[DEBUG] Intelligence Stats request from %s", r.RemoteAddr)
		stats, err := billingMgr.GetGlobalStats(r.Context())
		if err != nil {
			log.Printf("[ERROR] Failed to get global stats: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}))

	// --- Complete Manager Extensions (Phase 15) ---

	// Infrastructure: Proxmox Management
	mux.Handle("/infra/proxmox/containers", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var c domain.Container
			json.NewDecoder(r.Body).Decode(&c)
			if err := providerFactory.DeployContainer(r.Context(), &c); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "Container deployment initiated", "id": c.ID})
		} else {
			images, _ := providerFactory.ListImages(r.Context())
			json.NewEncoder(w).Encode(images)
		}
	}))

	// Infrastructure: OpenStack Orchestration
	mux.Handle("/infra/openstack/deploy", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Println("[Nebula] Triggering OpenStack Cloud-in-a-Box deployment via Blueprint...")
			// Logic to trigger OpenStack deployment
			json.NewEncoder(w).Encode(map[string]string{"message": "OpenStack deployment scheduled"})
		}
	}))

	// Networking: Multi-Provider Routing
	mux.Handle("/network/routes", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var input struct {
				Route          domain.Route `json:"route"`
				SourceProvider string       `json:"source_provider"`
				TargetProvider string       `json:"target_provider"`
			}
			json.NewDecoder(r.Body).Decode(&input)
			if err := networkMgr.CreateInterProviderRoute(r.Context(), &input.Route, input.SourceProvider, input.TargetProvider); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "Cross-provider route established"})
		}
	}))

	// Automated Infrastructure Provisioning (One-Click)
	mux.Handle("/infra/automated/provision", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var req domain.AutomaticProvisioningRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := autoProvisionUC.Execute(r.Context(), req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(map[string]string{"message": "Automated provisioning started", "preset": req.PresetID})
		} else {
			presets, _ := autoProvisionUC.ListPresets(r.Context())
			json.NewEncoder(w).Encode(presets)
		}
	}))

	// --- Open Cloud Vision: Regions & Zones (Phase 16) ---
	mux.Handle("/cloud/regions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var reg domain.Region
			json.NewDecoder(r.Body).Decode(&reg)
			if reg.ID == "" {
				reg.ID = domain.NewID()
			}
			regionRepo.Create(r.Context(), &reg)
			json.NewEncoder(w).Encode(reg)
		} else {
			regs, _ := regionRepo.List(r.Context())
			json.NewEncoder(w).Encode(regs)
		}
	}))

	mux.Handle("/cloud/zones", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var az domain.AvailabilityZone
			json.NewDecoder(r.Body).Decode(&az)
			if az.ID == "" {
				az.ID = domain.NewID()
			}
			zoneRepo.Create(r.Context(), &az)
			json.NewEncoder(w).Encode(az)
		} else {
			regionID := r.URL.Query().Get("region_id")
			if regionID != "" {
				azs, _ := zoneRepo.GetByRegion(r.Context(), regionID)
				json.NewEncoder(w).Encode(azs)
			} else {
				azs, _ := zoneRepo.List(r.Context())
				json.NewEncoder(w).Encode(azs)
			}
		}
	}))

	mux.Handle("/cloud/orchestrate/stack", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var input struct {
				ProjectID string   `json:"project_id"`
				Regions   []string `json:"regions"`
			}
			json.NewDecoder(r.Body).Decode(&input)
			if err := orchestrator.ProvisionMultiZoneStack(r.Context(), input.ProjectID, input.Regions); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "Global multi-region stack orchestration started"})
		}
	}))

	mux.Handle("/api/billing", authMiddleware.Authenticate(http.HandlerFunc(billingHandler.GetReport)))
	mux.Handle("/api/stats", authMiddleware.Authenticate(http.HandlerFunc(billingHandler.GetStats)))

	// Enterprise Hierarchy & Bare Metal
	mux.Handle("/api/organizations", authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			hierarchyHandler.CreateOrganization(w, r)
		default:
			hierarchyHandler.ListOrganizations(w, r)
		}
	})))
	mux.Handle("/api/departments", authMiddleware.Authenticate(http.HandlerFunc(hierarchyHandler.ListDepartments)))
	mux.Handle("/api/baremetal/nodes", authMiddleware.Authenticate(http.HandlerFunc(bmHandler.ListNodes)))
	mux.Handle("/api/baremetal/register", authMiddleware.Authenticate(http.HandlerFunc(bmHandler.RegisterNode)))
	mux.Handle("/api/baremetal/provision", authMiddleware.Authenticate(http.HandlerFunc(bmHandler.ProvisionNode)))
	mux.Handle("/api/baremetal/logs", authMiddleware.Authenticate(http.HandlerFunc(bmHandler.GetNodeLogs)))
	// Provider Management (Infrastructure Onboarding)
	mux.Handle("/api/providers", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			providers, _ := providerRepo.List(r.Context())
			json.NewEncoder(w).Encode(providers)
		case http.MethodPost:
			var p domain.Provider
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if p.ID == "" {
				p.ID = domain.NewID()
			}
			p.CreatedAt = time.Now()
			p.Status = "active"
			providerRepo.Create(r.Context(), &p)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(p)
		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			if err := providerRepo.Delete(r.Context(), id); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	handlerWithMetrics := metricsMiddleware.Metrics(mux)

	// Apply CORS as the absolute outermost layer
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", getEnv("CORS_ALLOWED_ORIGIN", "*"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handlerWithMetrics.ServeHTTP(w, r)
	})

	port := getEnv("PORT", "8000")
	host := getEnv("HOST", "0.0.0.0")
	addr := host + ":" + port
	fmt.Printf("NebulaOS API listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func seedEnterpriseDefaults(
	tenantRepo domain.TenantRepository,
	projectRepo domain.ProjectRepository,
	orgRepo domain.OrganizationRepository,
	deptRepo domain.DepartmentRepository,
	userRepo domain.UserRepository,
	resourceRepo domain.ResourceRepository,
	volumeRepo domain.VolumeRepository,
	bucketRepo domain.BucketRepository,
) {
	ctx := context.Background()

	tenantID := getEnv("SEED_TENANT_ID", "v-t1")
	orgID := getEnv("SEED_ORG_ID", "org-nebula-main")
	deptID := getEnv("SEED_DEPT_ID", "dept-core-infra")
	projectID := getEnv("SEED_PROJECT_ID", "v-p1")
	adminUser := getEnv("SEED_ADMIN_USERNAME", "admin")
	adminPass := getEnv("SEED_ADMIN_PASSWORD", "admin")
	adminEmail := getEnv("SEED_ADMIN_EMAIL", "admin@nebula.local")
	adminID := getEnv("SEED_ADMIN_ID", "u-admin")
	tenantName := getEnv("SEED_TENANT_NAME", "Nebula Global Corp")
	orgName := getEnv("SEED_ORG_NAME", "Nebula Global Corp")
	deptName := getEnv("SEED_DEPT_NAME", "Core Infrastructure")
	projectName := getEnv("SEED_PROJECT_NAME", "Nebula Core Project")

	// 1. Seed Default Tenant
	_, err := tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		fmt.Println("[Nebula] Seeding default tenant...")
		if err := tenantRepo.Create(ctx, &domain.Tenant{ID: tenantID, Name: tenantName, CreatedAt: time.Now()}); err != nil {
			log.Printf("Warning: failed to seed tenant: %v", err)
		}
	}

	// 2. Seed Default Organization
	_, err = orgRepo.GetByID(ctx, orgID)
	if err != nil {
		fmt.Println("[Nebula] Seeding default organization...")
		if err := orgRepo.Create(ctx, &domain.Organization{
			ID:        orgID,
			Name:      orgName,
			CreatedAt: time.Now(),
		}); err != nil {
			log.Printf("Warning: failed to seed organization: %v", err)
		}
	}

	// 3. Seed Default Department
	_, err = deptRepo.GetByID(ctx, deptID)
	if err != nil {
		fmt.Println("[Nebula] Seeding core infrastructure department...")
		if err := deptRepo.Create(ctx, &domain.Department{
			ID:             deptID,
			OrganizationID: orgID,
			Name:           deptName,
			CreatedAt:      time.Now(),
		}); err != nil {
			log.Printf("Warning: failed to seed department: %v", err)
		}
	}

	// 4. Seed Default Project
	_, err = projectRepo.GetByID(ctx, projectID)
	if err != nil {
		fmt.Println("[Nebula] Seeding default project...")
		if err := projectRepo.Create(ctx, &domain.Project{ID: projectID, TenantID: tenantID, Name: projectName, CreatedAt: time.Now()}); err != nil {
			log.Printf("Warning: failed to seed project: %v", err)
		}
	}

	// 5. Seed Admin User
	admin, err := userRepo.GetByUsername(ctx, adminUser)
	if err != nil || admin == nil {
		fmt.Println("[Nebula] Seeding default administrator...")
		tempAuth := services.NewInternalIdentityManager(userRepo, "temp")
		hash, hashErr := tempAuth.HashPassword(adminPass)
		if hashErr != nil {
			log.Printf("Warning: failed to hash admin password: %v", hashErr)
		} else {
			if err := userRepo.Create(ctx, &domain.User{
				ID:                 adminID,
				Username:           adminUser,
				PasswordHash:       hash,
				Email:              adminEmail,
				TenantID:           tenantID,
				MustChangePassword: true,
				CreatedAt:          time.Now(),
			}); err != nil {
				log.Printf("Warning: failed to seed admin user: %v", err)
			}
		}
	}

	// 6. Seed Default Resources if project is empty
	resources, err := resourceRepo.GetByProject(ctx, projectID)
	if err != nil {
		log.Printf("Warning: failed to check existing resources: %v", err)
	}
	if len(resources) == 0 {
		seedResourceID1 := getEnv("SEED_RESOURCE_NODE_ID", "v-n1")
		seedResourceID2 := getEnv("SEED_RESOURCE_VOL_ID", "v-s1")
		seedVolumeID := getEnv("SEED_VOLUME_ID", "v-vol-1")
		seedBucketID := getEnv("SEED_BUCKET_ID", "v-b1")
		seedBucketRegion := getEnv("SEED_BUCKET_REGION", "us-east-1")

		fmt.Println("[Nebula] Seeding initial resources...")
		if err := resourceRepo.Create(ctx, &domain.Resource{ID: seedResourceID1, ProjectID: projectID, Name: "Node-1", Type: domain.ComputeResource, State: "active", CreatedAt: time.Now()}); err != nil {
			log.Printf("Warning: failed to seed resource: %v", err)
		}
		if err := resourceRepo.Create(ctx, &domain.Resource{ID: seedResourceID2, ProjectID: projectID, Name: "Vol-1", Type: domain.StorageResource, State: "active", CreatedAt: time.Now()}); err != nil {
			log.Printf("Warning: failed to seed resource: %v", err)
		}
		if err := volumeRepo.Create(ctx, &domain.Volume{ID: seedVolumeID, ProjectID: projectID, Name: "Primary-OS-Disk", SizeGB: 100, State: "active", CreatedAt: time.Now()}); err != nil {
			log.Printf("Warning: failed to seed volume: %v", err)
		}
		if err := bucketRepo.Create(ctx, &domain.Bucket{ID: seedBucketID, ProjectID: projectID, Name: "Global-Assets", Region: seedBucketRegion, State: "active", CreatedAt: time.Now()}); err != nil {
			log.Printf("Warning: failed to seed bucket: %v", err)
		}
	}
}

func infrastructureSelfHeal() {
	fmt.Println("[Nebula] Checking infrastructure health...")
	// Minimal check: if we are in local environment, try to see if docker-compose is needed
	if _, err := os.Stat("deploy/local/docker-compose.yml"); err == nil {
		fmt.Println("[Nebula] Local infrastructure definition found. Ensuring services are up...")
		// Executing docker-compose up -d
		cmd := exec.Command("docker-compose", "-f", "deploy/local/docker-compose.yml", "up", "-d")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("[Nebula] Warning: Failed to auto-start infrastructure: %v\n", err)
		} else {
			fmt.Println("[Nebula] Infrastructure self-healing sequence completed.")
		}
	}
}
