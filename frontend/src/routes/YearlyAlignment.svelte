<script lang="ts">
    import {Check} from "lucide-svelte";
    import {DoYearlyAlignment, GetActiveMonthsNumber} from "../../wailsjs/go/main/App";
    import {onMount} from "svelte";
    import {MONEY_DIVIDER} from "../constants";
    import {fade} from 'svelte/transition';
    import {main} from "../../wailsjs/go/models";
    import {numberWithSpaces} from "../common_functions";
    import YearlyAlignmentResult = main.YearlyAlignmentResult;

    let alignmentResult = new YearlyAlignmentResult;

    let yearlyAlignmentForm = {
        Year: new Date().getFullYear(),
        ActiveMonths: 12, // Recalculated when the year is changed
    }

    function fetchActiveMonthsNumber(year: number): void {
        GetActiveMonthsNumber(year).then((result: number) => (yearlyAlignmentForm.ActiveMonths = result));
    }

    function displayYearlyAlignmentForm(): void {
        DoYearlyAlignment(yearlyAlignmentForm.Year).then((result: YearlyAlignmentResult) => (alignmentResult = result))
    }

    function getPaymentClass(paidSum, calculatedSum): string {
        return calculatedSum - paidSum > 0 ? "payment-required" : "nothing-to-pay";
    }

    function getPaymentText(paidSum, calculatedSum, sumString): string {
        if (calculatedSum - paidSum > 0) {
            return sumString + " за доплащане"
        } else if (calculatedSum - paidSum < 0) {
            return sumString + " надвнесени"
        }

        return "изравнени"
    }

    $: fetchActiveMonthsNumber(yearlyAlignmentForm.Year);

    onMount(() => {
        fetchActiveMonthsNumber(yearlyAlignmentForm.Year)
    });
</script>

<div class="input-box" id="input-box">

    <h2>Годишно изравняване</h2>

    <span>Година</span> <input type="number" id="yearly-alignment-year" class="year-input"
                               bind:value={yearlyAlignmentForm.Year}/>

    <span style="margin-left: 5px;">Активни месеци: {yearlyAlignmentForm.ActiveMonths}</span>

    <button class="btn btn-small" on:click={displayYearlyAlignmentForm}>
        <span><Check color="#444" size="20"/></span>
    </button>

    {#if alignmentResult.IsCalculated}
        <div id="yearly-alignment-result-block" in:fade={{duration:300}}>
            <table class="result-table">
                <tbody>
                <tr>
                    <td>Брутен годишен доход:</td>
                    <td><b>{numberWithSpaces(alignmentResult.YearlyGrossIncomeCents / MONEY_DIVIDER)} EUR</b></td>
                    <td style="width: 50px;"></td>
                    <td></td>
                    <td></td>
                </tr>
                <tr>
                    <td>Облагаем доход (-{alignmentResult.ExpensesPercentage}% разходи):</td>
                    <td><b>{numberWithSpaces(alignmentResult.TaxedYearlyIncomeCents / MONEY_DIVIDER)} EUR</b></td>
                    <td></td>
                    <td></td>
                    <td></td>
                </tr>
                <tr>
                    <td>Платени данъци:</td>
                    <td><b>{numberWithSpaces(alignmentResult.TaxesReallyPaidCents / MONEY_DIVIDER)} EUR</b></td>
                    <td></td>
                    <td>Изравнени данъци:</td>
                    <td>
                        <b class={getPaymentClass(alignmentResult.TaxesReallyPaidCents, alignmentResult.RecalculatedTaxCents)}>
                            {numberWithSpaces(alignmentResult.RecalculatedTaxCents / MONEY_DIVIDER)} EUR
                        </b>
                    </td>
                </tr>
                <tr>
                    <td>Платени осигуровки:</td>
                    <td><b>{numberWithSpaces(alignmentResult.SocialSecurityReallyPaidCents / MONEY_DIVIDER)} EUR</b>
                    </td>
                    <td></td>
                    <td>Изравнени осигуровки:</td>
                    <td>
                        <b class={getPaymentClass(alignmentResult.SocialSecurityReallyPaidCents, alignmentResult.RecalculatedSocialSecurityCents)}>
                            {numberWithSpaces(alignmentResult.RecalculatedSocialSecurityCents / MONEY_DIVIDER)} EUR
                        </b>
                    </td>
                </tr>
                </tbody>
            </table>

            <div class="alert alert-info" style="margin: 10px 0;">
                <p>Проверете тези суми при попълване на годишната данъчна декларация.</p>
                <p>Осигуровки: <b class={getPaymentClass(alignmentResult.SocialSecurityReallyPaidCents, alignmentResult.RecalculatedSocialSecurityCents)}>
                    {getPaymentText(
                        alignmentResult.SocialSecurityReallyPaidCents,
                        alignmentResult.RecalculatedSocialSecurityCents,
                        numberWithSpaces(alignmentResult.SocialSecurityDiffCents / MONEY_DIVIDER) + " EUR")
                    }
                </b></p>
                <p>Данъци: <b class={getPaymentClass(alignmentResult.TaxesReallyPaidCents, alignmentResult.RecalculatedTaxCents)}>
                    {getPaymentText(
                        alignmentResult.TaxesReallyPaidCents,
                        alignmentResult.RecalculatedTaxCents,
                        numberWithSpaces(alignmentResult.TaxesDiffCents / MONEY_DIVIDER) + " EUR")
                    }
                </b></p>

                <p>Разликата се доплаща до 30 април. Надвнесеното се приспада от бъдещите задължения или се възстановява:
                    търсете "Възстановяване на надвнесени суми" в портала на НАП.</p>

            </div>

            <table class="table">
                <thead>
                <tr>
                    <th>Месец</th>
                    <th>Брутен доход, EUR</th>
                    <th>Среден облагаем доход, EUR</th>
                    <th>Изравнен осигурителен доход, EUR</th>
                    <th>Изравнени осигуровки, EUR</th>
                    <th>Платени осигуровки, EUR</th>
                </tr>
                </thead>
                <tbody>
                {#each alignmentResult.Months as month}
                    <tr>
                        <td>{month.Month}</td>
                        <td>{numberWithSpaces(month.GrossIncomeCents / MONEY_DIVIDER)}</td>
                        <td>{numberWithSpaces(month.AverageTaxedIncomeCents / MONEY_DIVIDER)}</td>
                        <td>{numberWithSpaces(month.FinalSocSecIncomeCents / MONEY_DIVIDER)}</td>
                        <td>{numberWithSpaces(month.RecalculatedSocialSecurityCents / MONEY_DIVIDER)}</td>
                        <td>{numberWithSpaces(month.PaidSocialSecurityCents / MONEY_DIVIDER)}</td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

<style>
    .result-table {
        width: auto;
        margin-top: 15px;
        border-collapse: collapse;
    }

    .result-table td {
        padding: 8px 5px;
        text-align: left;
        border-bottom: 1px solid #c9c1ad;
    }

    .nothing-to-pay {
        color: #2a5007;
    }

    .payment-required {
        color: #700404;
    }
</style>