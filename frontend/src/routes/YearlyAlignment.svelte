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
        <div id="yearly-alignment-result-block" in:fade={{duration:300}} class="alert alert-info">
            <table class="result-table">
                <tbody>
                <tr>
                    <td>Брутен годишен доход:</td>
                    <td><b>{numberWithSpaces(alignmentResult.YearlyGrossIncomeCents / MONEY_DIVIDER)} EUR</b></td>
                </tr>
                <tr>
                    <td>Облагаем доход:</td>
                    <td><b>{numberWithSpaces(alignmentResult.TaxedYearlyIncomeCents / MONEY_DIVIDER)} EUR</b></td>
                </tr>
                <tr>
                    <td>Платени данъци:</td>
                    <td><b>{numberWithSpaces(alignmentResult.TaxesReallyPaidCents / MONEY_DIVIDER)} EUR</b></td>
                </tr>
                </tbody>
            </table>
            <table class="months-table">
                <thead>
                <tr>
                    <th>Месец</th>
                    <th>Брутен доход</th>
                    <th>Среден облагаем доход</th>
                    <th>Изравнен осигурителен доход</th>
                </tr>
                </thead>
                <tbody>
                {#each alignmentResult.Months as month}
                    <tr>
                        <td>{month.Month}</td>
                        <td>{numberWithSpaces(month.GrossIncomeCents / MONEY_DIVIDER)} EUR</td>
                        <td>{numberWithSpaces(month.AverageTaxedIncomeCents / MONEY_DIVIDER)} EUR</td>
                        <td>{numberWithSpaces(month.FinalInsuranceIncomeCents / MONEY_DIVIDER)} EUR</td>
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
        border-bottom: 1px solid #ddd;
    }

    .months-table {
        margin-top: 20px;
        width: 100%;
        border-collapse: collapse;
    }

    .months-table th,
    .months-table td {
        padding: 8px;
        text-align: left;
        border-bottom: 1px solid #ddd;
    }

    .months-table th {
        font-weight: bold;
        background-color: #f5f5f5;
    }
</style>