import { useEffect, useMemo, useState } from "react"

import styles from "./RelativeTime.module.css"

interface RelativeTimeProps {
  instant: Temporal.Instant
}

function usePeriodic<T>(f: () => T, interval: Temporal.Duration) {
  const [value, setValue] = useState(f())
  const intervalMs = interval.total({ unit: "milliseconds" })
  useEffect(() => {
    const timer = setInterval(() => setValue(f()), intervalMs)
    return () => clearInterval(timer)
  })
  return value
}

const CALENDAR_ISO8601 = Temporal.Calendar.from("iso8601")

function formatWithPrecision(fmt: Intl.RelativeTimeFormat, duration: Temporal.Duration) {
  duration = duration.round({
    largestUnit: "year",
    smallestUnit: "minute",
    relativeTo: Temporal.now.zonedDateTime(CALENDAR_ISO8601),
  })
  if (Math.abs(duration.years) > 0) {
    return fmt.format(duration.years, "year")
  } else if (Math.abs(duration.months) > 0) {
    return fmt.format(duration.months, "month")
  } else if (Math.abs(duration.days) > 0) {
    return fmt.format(duration.days, "day")
  } else if (Math.abs(duration.hours) > 0) {
    return fmt.format(duration.hours, "hour")
  } else if (Math.abs(duration.minutes) > 0) {
    return fmt.format(duration.minutes, "minute")
  } else {
    return "Just now"
  }
}

export function RelativeTime(props: RelativeTimeProps) {
  const now = usePeriodic(() => Temporal.now.instant(), Temporal.Duration.from({ minutes: 1 }))
  const fmt = useMemo(() => new Intl.RelativeTimeFormat(["en-US"]), [])
  return (
    <time
      className={styles.time}
      dateTime={props.instant.toString()}
      title={props.instant.round({ smallestUnit: "second" }).toString()}>
      {formatWithPrecision(fmt, props.instant.since(now))}
    </time>
  )
}
