package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/jacksonwendel/nebulaos/src/api/domain"
	"github.com/jacksonwendel/nebulaos/src/api/internal/api"
	"github.com/jacksonwendel/nebulaos/src/api/internal/api/middleware"
	"github.com/jacksonwendel/nebulaos/src/api/internal/infrastructure"
	"github.com/jacksonwendel/nebulaos/src/api/internal/services"
	"github.com/jacksonwendel/nebulaos/src/api/internal/usecase"
	"github.com/jacksonwendel/nebulaos/src/providers"
	"github.com/jacksonwendel/nebulaos/src/providers/aws"
	"github.com/jacksonwendel/nebulaos/src/providers/baremetal"
	"github.com/jacksonwendel/nebulaos/src/providers/keycloak"
	"github.com/jacksonwendel/nebulaos/src/providers/kubernetes"
	"github.com/jacksonwendel/nebulaos/src/providers/mock"
	"github.com/jacksonwendel/nebulaos/src/providers/openstack"
	"github.com/jacksonwendel/nebulaos/src/providers/proxmox"
	"github.com/jacksonwendel/nebulaos/src/providers/storage"
	"github.com/jacksonwendel/nebulaos/src/providers/traefik"
	"github.com/jacksonwendel/nebulaos/src/providers/vault"
	_ "github.com/lib/pq" // NEW
	"github.com/nats-io/nats.go"
)

func main() {
	fmt.Println("Starting NebulaOS Cloud API (Phase 11: Production Hardened)...")

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

	proxmoxProvider := proxmox.NewProxmoxProvider("https://pve.nebula.local/api2/json", "token-uuid")
	providerFactory.Register("proxmox", proxmoxProvider)

	k8sProvider := kubernetes.NewKubernetesProvider("kubeconfig-data")
	providerFactory.Register("kubernetes", k8sProvider)

	osProvider, _ := openstack.NewOpenStackProvider(gophercloud.AuthOptions{
		IdentityEndpoint: "http://openstack:5000/v3",
		Username:         "admin",
		Password:         "password",
	})
	if osProvider != nil {
		providerFactory.Register("openstack", osProvider)
	}

	bmProvider := baremetal.NewBareMetalProvider("admin", "password")
	providerFactory.Register("baremetal", bmProvider)

	awsProvider, _ := aws.NewAWSProvider(context.Background(), "us-east-1", "http://localhost:4566")
	if awsProvider != nil {
		providerFactory.Register("aws", awsProvider)
	}

	traefikProvider := traefik.NewTraefikProvider("./configs/traefik")
	keycloakProvider := keycloak.NewKeycloakProvider(kcURL, "nebula-api")
	vaultProvider := vault.NewVaultProvider(vaultURL, vaultToken)

	// Middleware
	metricsMiddleware := middleware.NewMetricsMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(keycloakProvider)
	auditMiddleware := middleware.NewAuditMiddleware(nc)

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

	if dbURL != "" {
		fmt.Println("Initializing PostgreSQL persistence layer...")
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
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
	} else {
		fmt.Println("Initializing In-Memory repositories (Development Mode)...")
		tenantRepo = infrastructure.NewInMemoryTenantRepository()
		projectRepo = infrastructure.NewInMemoryProjectRepository()
		resourceRepo = infrastructure.NewInMemoryResourceRepository()
		quotaRepo = infrastructure.NewInMemoryQuotaRepository()
		volumeRepo = infrastructure.NewInMemoryVolumeRepository()
		bucketRepo = infrastructure.NewInMemoryBucketRepository()
		policyRepo = infrastructure.NewInMemorySovereigntyPolicyRepository()
		securityGroupRepo = infrastructure.NewInMemorySecurityGroupRepository()
		tfStateRepo = infrastructure.NewInMemoryTerraformStateRepository()
		blueprintRepo = infrastructure.NewInMemoryBlueprintRepository()
	}

	fmt.Printf("Repositories initialized (Persistence: %v)\n", dbURL != "")
	_ = securityGroupRepo // Explicitly use to avoid lint error until usecases are added
	_ = tfStateRepo       // Explicitly use to avoid lint error until usecases are added
	_ = blueprintRepo     // Explicitly use to avoid lint error until usecases are added

	// Services
	policyService := services.NewSovereignPolicyService(policyRepo)
	billingMgr := infrastructure.NewSovereignBillingManager(resourceRepo, volumeRepo, bucketRepo)

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
	getResourceUC := usecase.NewGetResourceUseCase(resourceRepo)
	listResourcesUC := usecase.NewListResourcesByProjectUseCase(resourceRepo)

	createVolumeUC := usecase.NewCreateVolumeUseCase(volumeRepo, storageProvider)
	createBucketUC := usecase.NewCreateBucketUseCase(bucketRepo, storageProvider)

	requestCertUC := usecase.NewRequestCertificateUseCase(traefikProvider)
	storeSecretUC := usecase.NewStoreSecretUseCase(vaultProvider)
	complianceUC := usecase.NewGetComplianceReportUseCase(resourceRepo, projectRepo, quotaRepo)

	// Dependency Injection - Handlers
	tenantHandler := api.NewTenantHandler(createTenantUC, getTenantUC, listTenantsUC)
	projectHandler := api.NewProjectHandler(createProjectUC, getProjectUC, listProjectsUC)
	resourceHandler := api.NewResourceHandler(createResourceUC, getResourceUC, listResourcesUC)
	storageHandler := api.NewStorageHandler(createVolumeUC, createBucketUC)
	complianceHandler := api.NewComplianceHandler(complianceUC)
	billingHandler := api.NewBillingHandler(billingMgr)
	policyHandler := api.NewPolicyHandler(policyService)

	// Routes
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]string{
			"api":      "healthy",
			"nats":     "connected",
			"keycloak": "active",
			"vault":    "active",
		}
		if nc == nil || !nc.IsConnected() {
			status["nats"] = "disconnected"
		}
		w.Header().Set("Content-Type", "application/json")
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
	mux.Handle("/security-groups", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))))

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

	mux.Handle("/resources", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			resourceHandler.CreateResource(w, r)
		} else {
			resourceHandler.ListResources(w, r)
		}
	}))))

	mux.Handle("/storage/volumes", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(storageHandler.CreateVolume))))
	mux.Handle("/storage/buckets", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(storageHandler.CreateBucket))))

	// Billing & Governance
	mux.Handle("/billing/report", auditMiddleware.Audit(authMiddleware.Authenticate(http.HandlerFunc(billingHandler.GetReport))))
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
	mux.Handle("/marketplace/blueprints", authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})))

	// Apply metrics to ALL routes
	handlerWithMetrics := metricsMiddleware.Metrics(mux)

	port := getEnv("PORT", "8000")
	fmt.Printf("NebulaOS API listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithMetrics))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
