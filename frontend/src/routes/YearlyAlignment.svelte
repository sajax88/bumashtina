<script lang="ts">
    import {Check} from "lucide-svelte";
    import {DoYearlyAlignment, GetActiveMonthsNumber, LoadTaxesConfig} from "../../wailsjs/go/main/App";
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
        <div id="yearly-alignment-result-block" in:fade={{duration:300}} >
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
                    <td>Облагаем доход (-разходи):</td>
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
                        <b>{numberWithSpaces(alignmentResult.RecalculatedTaxCents / MONEY_DIVIDER)} EUR</b>
                        <!-- TODO: diff with paid -->
                    </td>
                </tr>
                <tr>
                    <td>Платени осигуровки:</td>
                    <td><b>{numberWithSpaces(alignmentResult.SocialSecurityReallyPaidCents / MONEY_DIVIDER)} EUR</b></td>
                    <td></td>
                    <td>Изравнени осигуровки</td>
                    <td>
                        <b>{numberWithSpaces(alignmentResult.RecalculatedSocialSecurityCents / MONEY_DIVIDER)} EUR</b>
                        <!-- TODO: diff with paid -->
                    </td>
                </tr>
                </tbody>
            </table>

            <!--
		// TODO
		//	След като знаеш окончателния си осигурителен доход,
		//		изчисляваш годишните осигуровки върху него. НАП сравнява тази сума с осигуровките,
		//		които вече си платил авансово:
		//	Платил си по-малко от дължимото → доплащаш разликата до 30.04
		//	Платил си повече от дължимото → надвнесеното се приспада от бъдещи задължения или ти се възстановява

		-->

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
                    <!-- TODO: beautify, maybe add colors where there're diffs -->
                    <tr>
                        <td>{month.Month}</td>
                        <td>{numberWithSpaces(month.GrossIncomeCents / MONEY_DIVIDER)}</td>
                        <td>{numberWithSpaces(month.AverageTaxedIncomeCents / MONEY_DIVIDER)}</td>
                        <td>{numberWithSpaces(month.FinalInsuranceIncomeCents / MONEY_DIVIDER)}</td>
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
        width:auto;
        margin-top: 15px;
        border-collapse: collapse;
    }

    .result-table td {
        padding: 8px 5px;
        text-align: left;
        border-bottom: 1px solid #c9c1ad;
    }
</style>