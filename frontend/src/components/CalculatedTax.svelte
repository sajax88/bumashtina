<script lang="ts">
    import {MONEY_DIVIDER} from "../constants";
    import {Info} from "lucide-svelte";
    import {fade} from 'svelte/transition';

    export let taxCalculationResult;

    // Thank you, dude from stackoverflow
    function numberWithSpaces(x) {
        let parts = x.toString().split(".");
        parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, " ");
        return parts.join(".");
    }
</script>

<div in:fade={{duration:300}} class="alert alert-info" id="tax-calculator-result">
    <table class="table" id="tax-calculator-result-table">
        <tbody>
        <tr>
            <td><small>доход за месеци {taxCalculationResult.MonthStart}-{taxCalculationResult.MonthEnd}</small></td>
            <td>{numberWithSpaces(taxCalculationResult.TotalIncomeCents / MONEY_DIVIDER)}</td>
        </tr>
        <tr>
            <td><small>- {taxCalculationResult.TaxesConfig.ExpensesPercentage}% разходи</small></td>
            <td>-{numberWithSpaces(taxCalculationResult.ExpensesCents / MONEY_DIVIDER)}</td>
        </tr>
        <tr>
            <td><small>- осигуровки</small></td>
            <td>-{numberWithSpaces(taxCalculationResult.PaidInsuranceCents / MONEY_DIVIDER)}</td>
        </tr>
        <tr>
            <td><small>* {taxCalculationResult.TaxesConfig.TaxPercentage}%</small></td>
            <td>= <b>{numberWithSpaces(taxCalculationResult.TaxCents / MONEY_DIVIDER)} EUR</b></td>
        </tr>
        </tbody>
    </table>
    {#if taxCalculationResult.Notes}
        <div class="alert-small-note">
            <small>
                <span><Info color="#79b4d1" size="18"/>{taxCalculationResult.Notes}</span>
            </small>
        </div>
    {/if}
</div>

<style>
    #tax-calculator-result-table  {
        margin-top: 0;
    }

    #tax-calculator-result-table td {
        border-left: 0;
        border-right: 0;
        text-align: right;
    }

    #tax-calculator-result .alert-small-note {
        margin-top: 10px;
    }
</style>