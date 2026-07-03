import { MONEY_DIVIDER } from "./constants";

export function moneyToCents(sum: string): number
{
    return Math.round(parseFloat(sum) * MONEY_DIVIDER)
}