<script lang="ts">
    import {MONEY_DIVIDER} from "../constants";
    import {Info} from "lucide-svelte";
    import {numberWithSpaces} from "../common_functions";

    export let dataSingle;
</script>
<div id="tax-calculation-for-month">
        <table class="table" id="tax-calculation-for-month-table">
            <tbody>
            <tr>
                <td>доход за месец</td>
                <td>{numberWithSpaces(dataSingle.MonthIncomeCents / MONEY_DIVIDER)}</td>
            </tr>
            <tr>
                <td>- {dataSingle.TaxesConfig.ExpensesPercentage}% разходи</td>
                <td>
                    -{numberWithSpaces(dataSingle.ExpensesCents / MONEY_DIVIDER)}
                </td>
            </tr>
            <tr>
                <td>- осигуровки</td>
                <td>
                    {#if dataSingle.SocialSecurityReallyPaidCents}
                        -{numberWithSpaces(dataSingle.SocialSecurityReallyPaidCents / MONEY_DIVIDER)}
                    {:else}
                        -{numberWithSpaces(dataSingle.SocialSecurityToPayCents / MONEY_DIVIDER)}
                    {/if}
                </td>
            </tr>
            <tr>
                <td>* {dataSingle.TaxesConfig.TaxPercentage}%</td>
                <td>= <b>{numberWithSpaces(dataSingle.TaxesToPayCents / MONEY_DIVIDER)} EUR</b></td>
            </tr>
            </tbody>
        </table>
        {#if !dataSingle.SocialSecurityReallyPaidCents}
            <div class="alert-small-note">
                <small>
                    <span><Info color="#79b4d1" size="18" style="margin-right: 5px" /> Въз основа на изчислените осигуровки</span>
                </small>
            </div>
        {/if}
</div>
<style>
    #tax-calculation-for-month-table {
        margin-top: 0;
        font-size:small;
        width: auto;
    }

    #tax-calculation-for-month-table td {
        border:0;
        text-align: right;
        padding: 4px 10px;
    }

    #tax-calculation-for-month .alert-small-note {
        margin-top: 10px;
    }
</style>