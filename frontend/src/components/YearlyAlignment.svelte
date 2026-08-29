<script lang="ts">
    import {Check, Scale} from "lucide-svelte";
    import {GetActiveMonthsNumber, DoYearlyAlignment} from "../../wailsjs/go/main/App";
    import {onMount} from "svelte";
    import {MONEY_DIVIDER} from "../constants";
    import {fade} from 'svelte/transition';
    import {main} from "../../wailsjs/go/models";
    import YearlyAlignmentResult = main.YearlyAlignmentResult;
    import {numberWithSpaces} from "../common_functions";

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

<div id="yearly-alignment-box">
    <div class="form-row">
        <div class="form-group">
            <button class="btn"
                    on:click={() => {document.getElementById('yearly-alignment-block').style.display = 'block';}}>
                <span><Scale color="#444" size="20"/> Годишно изравняване</span>
            </button>

            <div id="yearly-alignment-block" class="hidden-form-block" style="display: none;">
                <small>Година</small> <input type="number" id="yearly-alignment-year" class="year-input"
                                             bind:value={yearlyAlignmentForm.Year}/>

                <small style="margin-left: 5px;">Активни месеци: {yearlyAlignmentForm.ActiveMonths}</small>

                <button class="btn btn-small" on:click={displayYearlyAlignmentForm}>
                    <span><Check color="#444" size="20"/></span>
                </button>

                {#if alignmentResult.IsCalculated}
                    <div id="yearly-alignment-result-block"  in:fade={{duration:300}} class="alert alert-info">
                        Брутен годишен доход: <b>{numberWithSpaces(alignmentResult.YearlyGrossIncomeCents / MONEY_DIVIDER)} EUR</b><br>
                        Облагаем доход: <b>{numberWithSpaces(alignmentResult.TaxedYearlyIncomeCents / MONEY_DIVIDER)} EUR</b><br>
                        Платени данъци: <b>{numberWithSpaces(alignmentResult.TaxesReallyPaidCents / MONEY_DIVIDER)} EUR</b><br>
                        <!-- TODO
                        <table class="months-table">
                            <thead>
                            <tr>
                                <th>Месец</th>
                                <th>Брутен доход</th>
                                <th>Изравнен облагаем доход</th>
                                <th>Платени данъци</th>
                            </tr>
                            </thead>
                            <tbody>
                            {#each alignmentResult.Months as month}
                                <tr>
                                    <td>{month.MonthNumber}</td>
                                    <td>{numberWithSpaces(month.GrossIncomeCents / MONEY_DIVIDER)} EUR</td>
                                    <td>{numberWithSpaces(month.TaxedIncomeCents / MONEY_DIVIDER)} EUR</td>
                                    <td>{numberWithSpaces(month.TaxesPaidCents / MONEY_DIVIDER)} EUR</td>
                                </tr>
                            {/each}
                            </tbody>
                        </table>-->
                    </div>
                {/if}
            </div>
        </div>
    </div>
</div>

<style>
    #yearly-alignment-box {
        padding-top: 20px;
    }

    #yearly-alignment-block {
        padding-top: 10px;
    }

    .months-table {
        margin-top: 15px;
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