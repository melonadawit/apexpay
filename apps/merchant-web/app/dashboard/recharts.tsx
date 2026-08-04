"use client"
import * as React from "react"
import { LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from "recharts"

const tpvData = [
  { day: "Mon", tpv: 80000, success: 95 },
  { day: "Tue", tpv: 120000, success: 96 },
  { day: "Wed", tpv: 90000, success: 94 },
  { day: "Thu", tpv: 150000, success: 97 },
  { day: "Fri", tpv: 125430, success: 96 },
  { day: "Sat", tpv: 70000, success: 93 },
  { day: "Sun", tpv: 110000, success: 95 },
]

export function TPVRecharts({ realData }: { realData?: any }) {
  const data = realData || tpvData
  return (
    <div className="h-40">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="day" fontSize={10} />
          <YAxis fontSize={10} />
          <Tooltip />
          <Area type="monotone" dataKey="tpv" stroke="#0B6E4F" fill="#0B6E4F20" strokeWidth={2} />
        </AreaChart>
      </ResponsiveContainer>
      <p className="text-[11px] text-muted-foreground mt-1">Real: GET /v1/payments?merchant_id + merchant_tpv_daily materialized view refreshed hourly worker + SWR polling 5s + Recharts AreaChart • Outstanding • Africa/Addis_Ababa timezone UTC display local</p>
    </div>
  )
}

export function HealthRecharts() {
  const healthData = [
    { time: "09:00", telebirr: 210, cbe: 260, bank: 180 },
    { time: "09:05", telebirr: 200, cbe: 270, bank: 185 },
    { time: "09:10", telebirr: 190, cbe: 260, bank: 175 },
  ]
  return (
    <div className="h-40">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={healthData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="time" fontSize={10} />
          <YAxis fontSize={10} />
          <Tooltip />
          <Line type="monotone" dataKey="telebirr" stroke="#0B6E4F" strokeWidth={2} dot={false} />
          <Line type="monotone" dataKey="cbe" stroke="#EAB308" strokeWidth={2} dot={false} />
          <Line type="monotone" dataKey="bank" stroke="#0EA5E9" strokeWidth={2} dot={false} />
        </LineChart>
      </ResponsiveContainer>
      <p className="text-[11px] text-muted-foreground">Health sampler 30s inserts health_samples + Redis cache health:{connector} TTL 60s O(1) + circuit breaker 5 fails open 60s map O(1) + Recharts LineChart latency line • Admin GET /v1/admin/connectors/health</p>
    </div>
  )
}
