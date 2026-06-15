import React, { createContext, useContext, useState, useCallback } from 'react';
import en from '../locales/en';
import ptBR from '../locales/pt-BR';

type Locale = 'en' | 'pt-BR';
type Translations = typeof en;

const locales: Record<Locale, Translations> = { en, 'pt-BR': ptBR };

interface LocaleContextValue {
  locale: Locale;
  t: Translations;
  setLocale: (locale: Locale) => void;
}

const LocaleContext = createContext<LocaleContextValue>({
  locale: 'en',
  t: en,
  setLocale: () => {},
});

export function LocaleProvider({ children }: { children: React.ReactNode }) {
  const [locale, setLocale] = useState<Locale>(() => {
    const saved = localStorage.getItem('nebula_locale');
    if (saved === 'pt-BR' || saved === 'en') return saved;
    return navigator.language.startsWith('pt') ? 'pt-BR' : 'en';
  });

  const handleSetLocale = useCallback((newLocale: Locale) => {
    localStorage.setItem('nebula_locale', newLocale);
    setLocale(newLocale);
  }, []);

  return (
    <LocaleContext.Provider value={{ locale, t: locales[locale], setLocale: handleSetLocale }}>
      {children}
    </LocaleContext.Provider>
  );
}

export function useLocale() {
  return useContext(LocaleContext);
}
