"use client"
import * as React from "react"

type Language = "en" | "am"

interface LanguageContextType {
  language: Language
  setLanguage: (lang: Language) => void
  t: (en: string, am: string) => string
}

const LanguageContext = React.createContext<LanguageContextType>({
  language: "en",
  setLanguage: () => {},
  t: (en) => en,
})

export function LanguageProvider({ children }: { children: React.ReactNode }) {
  const [language, setLanguageState] = React.useState<Language>("en")

  React.useEffect(() => {
    const saved = localStorage.getItem("apexpay_lang") as Language
    if (saved && (saved === "en" || saved === "am")) {
      setLanguageState(saved)
    }
  }, [])

  const setLanguage = (lang: Language) => {
    setLanguageState(lang)
    localStorage.setItem("apexpay_lang", lang)
  }

  const t = React.useCallback(
    (en: string, am: string) => {
      return language === "en" ? en : am
    },
    [language]
  )

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  )
}

export const useLanguage = () => React.useContext(LanguageContext)
