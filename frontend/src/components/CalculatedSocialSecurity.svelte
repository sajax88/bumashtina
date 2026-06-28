<script lang="ts">
    import {MONEY_DIVIDER} from "../constants";
    import {numberWithSpaces} from "../common_functions";

    export let dataSingle;
    export let socialSecurityParts;
</script>

<table class="table social-security-table">
    <tr>
        <td>ДОО {dataSingle.TaxesConfig.PensionPercentagePartOne}% {#if dataSingle.Settings.IsPregnancyInsuranceEnabled} + ОЗМ {dataSingle.TaxesConfig.PregnancyInsurancePercentage}%{/if}</td>
        <td>
            {numberWithSpaces(socialSecurityParts.PensionPartOneCents / MONEY_DIVIDER)}
        </td>
    </tr>
    <tr>
        <td>ДЗПО {dataSingle.TaxesConfig.PensionPercentagePartTwo}% </td>
        <td>
            + {numberWithSpaces(socialSecurityParts.PensionPartTwoCents / MONEY_DIVIDER)}
        </td>
    </tr>
    <tr>
        <td>НЗОК {dataSingle.TaxesConfig.HealthInsurancePercentage}% </td>
        <td>
            + {numberWithSpaces(socialSecurityParts.HealthInsuranceCents / MONEY_DIVIDER)}
        </td>
    </tr>
    <tr>
        <td></td>
        <td>= <b>{numberWithSpaces(dataSingle.SocialSecurityToPayCents / MONEY_DIVIDER)} EUR</b></td>
    </tr>
</table>

<style>
    .social-security-table {
        margin-top: 0;
        font-size:small;
        width: auto;
        border-collapse: collapse;
    }

    .social-security-table td {
        border:0;
        text-align: right;
        padding: 4px 10px;
    }
</style>