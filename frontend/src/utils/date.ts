import { DateTime } from "luxon";

export function localizedDateFromISO(
    isoString: string,
    format = DateTime.DATE_FULL
): string {
    const locale = navigator.language || "en-US";
    return DateTime.fromISO(isoString).setLocale(locale).toLocaleString(format);
}
