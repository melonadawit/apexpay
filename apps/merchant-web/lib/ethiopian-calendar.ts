// Ethiopian Calendar Utility — Enkutatash (Ethiopian New Year) — Ethiopia Law Compliance
// Ethiopian calendar: 13 months — 12 months 30 days + Pagume 5/6 days
// Months: Meskerem (1), Tikimt (2), Hidar (3), Tahsas (4), Tir (5), Yekatit (6), Megabit (7), Miyazya (8), Ginbot (9), Sene (10), Hamle (11), Nehasse (12), Pagume (13)
// Enkutatash: Meskerem 1 = September 11 (or 12 in Gregorian leap year before Feb 29) per Ethiopian calendar
// For payroll: cutoff 25th Gregorian, but Ethiopian date for local compliance reports, payslip bilingual EN/AM
// O(1) conversion for payroll calendar view Recharts per Ethiopia business practice

export const ethiopianMonths = [
  "Meskerem", "Tikimt", "Hidar", "Tahsas", "Tir", "Yekatit",
  "Megabit", "Miyazya", "Ginbot", "Sene", "Hamle", "Nehasse", "Pagume"
]

export const ethiopianMonthsAmharic = [
  "መስከረም", "ጥቅምት", "ኅዳር", "ታኅሣሥ", "ጥር", "የካቲት",
  "መጋቢት", "ሚያዝያ", "ግንቦት", "ሰኔ", "ሐምሌ", "ነሐሴ", "ጷጉሜ"
]

// Simplified conversion: Ethiopian year = Gregorian year - 7 or -8 depending on month
// Approximate algorithm O(1) for payroll purposes — for exact conversion need full algorithm with leap years
// For outstanding UI, we use approximate + show both Gregorian and Ethiopian

export function gregorianToEthiopian(gregDate: Date): { year: number; month: number; day: number; monthName: string; monthNameAm: string; formatted: string; formattedAm: string } {
  // Simplified: Ethiopian New Year Enkutatash is Meskerem 1 = Sept 11 (Gregorian)
  // So: If Gregorian month >= Sept (9) and day >=11, Ethiopian year = Gregorian year -7, else -8
  // This is approximate O(1) for payroll calendar view
  const gYear = gregDate.getFullYear()
  const gMonth = gregDate.getMonth() + 1 // 1-12
  const gDay = gregDate.getDate()

  // Determine Ethiopian year
  let ethYear: number
  if (gMonth > 9 || (gMonth === 9 && gDay >= 11)) {
    ethYear = gYear - 7
  } else {
    ethYear = gYear - 8
  }

  // Determine Ethiopian month/day approximate
  // Enkutatash Meskerem 1 = Sept 11 Gregorian
  // So: Sept 11 = Meskerem 1, Oct 11 = Tikimt 1, Nov 10 = Hidar 1, etc. Approximate 30 days per month
  // For outstanding UI, we calculate days since last Enkutatash
  const enkutatash_this_year = new Date(gYear, 8, 11) // Sept 11 this Gregorian year (month 8 = Sept)
  let daysSinceEnkutatash: number
  let ethMonth: number
  let ethDay: number

  if (gregDate >= enkutatash_this_year) {
    // After Enkutatash this year
    daysSinceEnkutatash = Math.floor((gregDate.getTime() - enkutatash_this_year.getTime()) / (1000*60*60*24))
  } else {
    // Before Enkutatash this year, use previous year's Enkutatash
    const enkutatash_prev = new Date(gYear - 1, 8, 11)
    daysSinceEnkutatash = Math.floor((gregDate.getTime() - enkutatash_prev.getTime()) / (1000*60*60*24))
  }

  ethMonth = Math.floor(daysSinceEnkutatash / 30) + 1
  ethDay = (daysSinceEnkutatash % 30) + 1

  if (ethMonth > 13) {
    ethMonth = 13
    ethDay = Math.min(ethDay, 6) // Pagume 5/6 days
  }
  if (ethMonth < 1) ethMonth = 1
  if (ethDay < 1) ethDay = 1
  if (ethDay > 30 && ethMonth <= 12) ethDay = 30
  if (ethMonth === 13 && ethDay > 6) ethDay = 6

  const monthName = ethiopianMonths[ethMonth - 1] || "Meskerem"
  const monthNameAm = ethiopianMonthsAmharic[ethMonth - 1] || "መስከረም"

  return {
    year: ethYear,
    month: ethMonth,
    day: ethDay,
    monthName,
    monthNameAm,
    formatted: `${monthName} ${ethDay}, ${ethYear}`,
    formattedAm: `${monthNameAm} ${ethDay}, ${ethYear}`,
  }
}

export function formatEthiopianDate(gregDate: Date): string {
  const eth = gregorianToEthiopian(gregDate)
  return `${eth.formatted} (${eth.formattedAm})`
}

export function formatEthiopianDateShort(gregDate: Date): string {
  const eth = gregorianToEthiopian(gregDate)
  return `${eth.monthNameAm} ${eth.day} • ${eth.monthName} ${eth.day}`
}

// Ethiopian calendar months for payroll calendar view Recharts per Ethiopia business practice
export function getEthiopianMonthsForYear(gYear: number): Array<{ gregorianMonth: number; gregorianYear: number; ethiopianMonth: number; ethiopianYear: number; monthName: string; monthNameAm: string; cutoffDate: Date; disbursalDate: Date; payDate: Date }> {
  const result = []
  for (let m = 1; m <= 12; m++) {
    const gregDate = new Date(gYear, m - 1, 15)
    const eth = gregorianToEthiopian(gregDate)
    const cutoffDate = new Date(gYear, m - 1, 25) // cutoff 25th Ethiopia business practice
    const disbursalDate = new Date(gYear, m - 1, 30) // disbursal 30th
    const payDate = new Date(gYear, m, 0) // last day of month
    result.push({
      gregorianMonth: m,
      gregorianYear: gYear,
      ethiopianMonth: eth.month,
      ethiopianYear: eth.year,
      monthName: eth.monthName,
      monthNameAm: eth.monthNameAm,
      cutoffDate,
      disbursalDate,
      payDate,
    })
  }
  return result
}

// Enkutatash date for a given Gregorian year
export function getEnkutatashDate(gYear: number): Date {
  // Enkutatash is Meskerem 1 = Sept 11 Gregorian (Sept 12 in leap year before Feb 29? Simplified Sept 11)
  return new Date(gYear, 8, 11) // Month 8 = September (0-indexed)
}

// Check if date is Enkutatash (Ethiopian New Year)
export function isEnkutatash(gregDate: Date): boolean {
  const eth = gregorianToEthiopian(gregDate)
  return eth.month === 1 && eth.day === 1
}

// Ethiopian public holidays per Labour Proclamation 1156/2019 + Ethiopian calendar holidays
// For payroll OT calculation: holiday OT 2.0x per Art 90(2)
export const ethiopianPublicHolidays = [
  { month: 9, day: 11, name: "Enkutatash • Ethiopian New Year • አዲስ አመት", nameAm: "እንቁጣጣሽ • አዲስ አመት", type: "public_holiday", ot_rate: 2.0 },
  { month: 9, day: 27, name: "Meskel • Finding of True Cross • መስቀል", nameAm: "መስቀል", type: "public_holiday", ot_rate: 2.0 },
  { month: 1, day: 7, name: "Genna • Ethiopian Christmas • ገና", nameAm: "ገና", type: "public_holiday", ot_rate: 2.0 },
  { month: 1, day: 19, name: "Timket • Epiphany • ጥምቀት", nameAm: "ጥምቀት", type: "public_holiday", ot_rate: 2.0 },
  { month: 3, day: 2, name: "Adwa Victory Day • የዓድዋ ድል • 1896", nameAm: "የዓድዋ ድል", type: "public_holiday", ot_rate: 2.0 },
  { month: 5, day: 1, name: "Labour Day • የሰራተኞች ቀን", nameAm: "የሰራተኞች ቀን", type: "public_holiday", ot_rate: 2.0 },
  { month: 5, day: 5, name: "Patriots Victory Day • የአርበኞች ቀን", nameAm: "የአርበኞች ቀን", type: "public_holiday", ot_rate: 2.0 },
  { month: 5, day: 28, name: "Downfall of Derg • የደርግ ውድቀት", nameAm: "የደርግ ውድቀት", type: "public_holiday", ot_rate: 2.0 },
]

export function isPublicHoliday(gregDate: Date): { isHoliday: boolean; holiday?: typeof ethiopianPublicHolidays[0] } {
  const month = gregDate.getMonth() + 1
  const day = gregDate.getDate()
  const found = ethiopianPublicHolidays.find(h => h.month === month && h.day === day)
  return { isHoliday: !!found, holiday: found }
}
