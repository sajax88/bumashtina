<script lang="ts">
    import {
        GenerateDeclarationOne,
        LoadIncomeDataForMonth,
        LoadTaxesConfigLabels,
        UpdateForm
    } from "../../wailsjs/go/main/App.js";
    import {onMount} from 'svelte';
    import {Ban, BookText, Check, CircleArrowLeft, Edit, Eye, EyeOff, Save} from "lucide-svelte";
    import CalculatedTaxForMonth from "../components/CalculatedTaxForMonth.svelte";
    import {numberWithSpaces} from "../common_functions";
    import CalculatedSocialSecurity from "../components/CalculatedSocialSecurity.svelte";

    export let params;

    const MONEY_DIVIDER = 100;

    let dataSingle;
    let socialSecurityParts;
    let socialSecurityPaidParts;

    let month = parseInt(params.month);
    let year = parseInt(params.year);

    let taxesLabels;
    let showTaxesConfig = false;

    let paidInsuranceValues = {
        PensionPartOne: 0,
        PensionPartTwo: 0,
        HealthInsurance: 0,
    }

    onMount(async () => {
        await load_data_for_month()
    });

    async function load_data_for_month(): Promise<void> {
        LoadIncomeDataForMonth(month, year).then(function (result) {
            dataSingle = result
            socialSecurityParts = dataSingle.SocialSecurityToPayParts
            socialSecurityPaidParts = dataSingle.SocialSecurityReallyPaidParts
            if (dataSingle.SocialSecurityReallyPaidCents) {
                paidInsuranceValues.PensionPartOne = socialSecurityPaidParts.PensionPartOneCents / MONEY_DIVIDER
                paidInsuranceValues.PensionPartTwo = socialSecurityPaidParts.PensionPartTwoCents / MONEY_DIVIDER
                paidInsuranceValues.HealthInsurance = socialSecurityPaidParts.HealthInsuranceCents / MONEY_DIVIDER

            } else {
                // Set to calculated sums
                paidInsuranceValues.PensionPartOne = socialSecurityParts.PensionPartOneCents / MONEY_DIVIDER
                paidInsuranceValues.PensionPartTwo = socialSecurityParts.PensionPartTwoCents / MONEY_DIVIDER
                paidInsuranceValues.HealthInsurance = socialSecurityParts.HealthInsuranceCents / MONEY_DIVIDER
            }
        });
        LoadTaxesConfigLabels().then(function (result) {
            taxesLabels = result
        });
    }

    function decl1(): void {
        GenerateDeclarationOne(month, year);
    }

    function showPaidInsuranceEditInput() {
        document.getElementById(`paid-insurance-inputs`).style.display = "inline-block";
        document.getElementById(`paid-insurance-save-button`).style.display = "inline-block";

        document.getElementById(`paid-insurance-value`).style.display = "none";
    }

    async function savePaidInsurance() {
        dataSingle.SocialSecurityReallyPaidParts.PensionPartOneCents = moneyToCents(paidInsuranceValues.PensionPartOne);
        dataSingle.SocialSecurityReallyPaidParts.PensionPartTwoCents = moneyToCents(paidInsuranceValues.PensionPartTwo);
        dataSingle.SocialSecurityReallyPaidParts.HealthInsuranceCents = moneyToCents(paidInsuranceValues.HealthInsurance);
        await UpdateForm(dataSingle).then(async function (result) {
            if (result != "") {
                document.getElementById(`paid-insurance-inputs`).style.display = "none";
                document.getElementById(`paid-insurance-save-button`).style.display = "none";
                document.getElementById(`paid-insurance-value`).style.display = "inline";

                await load_data_for_month()
            }
        });
    }

    function moneyToCents(sum: string): number
    {
        return parseInt(Math.round(parseFloat(sum) * MONEY_DIVIDER))
    }
</script>

<main>
    <div class="input-box" id="input-box">
        <a href="#/data">
            <CircleArrowLeft color="#444" size="20"/>
        </a>

        <h2>Въведени данни за {month}/{year}</h2>
        {#if dataSingle}
            <div id="data-single-buttons-block">
                <button class="btn btn-large" id="declaration-one-button" on:click={() => decl1()}>
                    <span><BookText color="#444" size="20"/> Генерирай Декларация 1</span>
                </button>
            </div>

            <table class="table">
                <tbody>
                <tr>
                    <td>Доход за месец</td>
                    <td colspan="2">{numberWithSpaces(dataSingle.MonthIncomeCents / MONEY_DIVIDER)}</td>
                </tr>
                <tr>
                    <td>Осигурителен доход</td>
                    <td colspan="2">{numberWithSpaces(dataSingle.TaxedIncomeCents / MONEY_DIVIDER)}</td>
                </tr>
                <tr>
                    <td>Изчислен данък<br></td>
                    <td colspan="2">
                       <CalculatedTaxForMonth dataSingle={dataSingle}/>
                    </td>
                </tr>
                <tr>
                    <td>Изчислени осигуровки</td>
                    <td colspan="2">
                        <CalculatedSocialSecurity dataSingle={dataSingle} socialSecurityParts={socialSecurityParts} />
                    </td>
                </tr>
                <tr>
                    <td>Платени осигуровки</td>
                    <td>
                        <span id="paid-insurance-value">
                            {#if dataSingle.SocialSecurityReallyPaidCents}
                                {numberWithSpaces(socialSecurityPaidParts.PensionPartOneCents / MONEY_DIVIDER)} ДОО {#if dataSingle.Settings.IsPregnancyInsuranceEnabled} + ОЗМ{/if}<br>
                                +{numberWithSpaces(socialSecurityPaidParts.PensionPartTwoCents / MONEY_DIVIDER)} ДЗПО<br>
                                +{numberWithSpaces(socialSecurityPaidParts.HealthInsuranceCents / MONEY_DIVIDER)} НЗОК<br>
                                = <b>{numberWithSpaces(dataSingle.SocialSecurityReallyPaidCents / MONEY_DIVIDER)} EUR</b>
                           {:else}
                                -
                            {/if}
                        </span>

                        <span id="paid-insurance-inputs" style="display: none">
                        <label class="paid-insurance-label">ДОО {#if dataSingle.Settings.IsPregnancyInsuranceEnabled} + ОЗМ{/if}: </label><input
                                class="paid-insurance-input" type="text"
                                bind:value={paidInsuranceValues.PensionPartOne}
                        /><br/>
                        <label class="paid-insurance-label">ДЗПО: </label><input
                                class="paid-insurance-input" type="text"
                                bind:value={paidInsuranceValues.PensionPartTwo}
                        /><br/>
                        <label class="paid-insurance-label">НЗОК: </label><input
                                class="paid-insurance-input" type="text"
                                bind:value={paidInsuranceValues.HealthInsurance}
                        />
                            </span><br/>

                        <button style="display: none" class="btn btn-small" id="paid-insurance-save-button"
                                on:click="{() => savePaidInsurance()}">
                            <Save color="#444" size="20"/>
                        </button>
                    </td>
                    <td>
                        <button class="btn btn-small" on:click="{() => showPaidInsuranceEditInput()}">
                            <Edit color="#444" size="20"/>
                        </button>
                    </td>
                </tr>
                <tr>
                    <td>Работни дни</td>
                    <td colspan="2">{dataSingle.WorkDaysTotal}</td>
                </tr>
                {#if dataSingle.WorkDaysSickLeave}
                    <tr>
                        <td>В болничен</td>
                        <td colspan="2">{dataSingle.WorkDaysSickLeave}</td>
                    </tr>
                {/if}

                {#if dataSingle.DayStart}
                    <tr>
                        <td>Започване на дейност от ден</td>
                        <td colspan="2">{dataSingle.DayStart}</td>
                    </tr>
                {/if}

                {#if dataSingle.DayEnd}
                    <tr>
                        <td>Приключване на дейност от ден</td>
                        <td colspan="2">{dataSingle.DayEnd}</td>
                    </tr>
                {/if}
            </table>
            <br>
            <p>
                <small><i>
                    Изчисленията стават на базата на таксите и осигуровките, които са били активни по време на въвеждането на данните.
                    За да промените таксите и осигуровките за следващите месеци, отидете в <a href="#/settings">"Настройки"</a>
                </i></small>
            </p>
            <button class="btn btn-small" on:click={() => showTaxesConfig = !showTaxesConfig}>
                            <span>{#if showTaxesConfig}<EyeOff color="#444" size="16"/>{:else}<Eye color="#444"
                                                                                                   size="16"/>{/if}
                                Виж настройки, данъци и осигуровки за този месец</span>
            </button>
            {#if showTaxesConfig}
                <h3>Данъци и осигуровки за {month}/{year}</h3>
                <table class="table">
                    <tbody>
                    <tr>
                        <td>Осигуряване за общо заболяване и майчинство</td>
                        <td>
                            {#if dataSingle.Settings.IsPregnancyInsuranceEnabled}
                                <Check color="#7a9423" size="20"/>
                            {:else}
                                <Ban color="#bf550f" size="20"/>
                            {/if}
                        </td>
                    </tr>
                    {#each Object.entries(dataSingle.TaxesConfig) as [key, value]}
                        <tr>
                            <td>{taxesLabels[Object.keys(dataSingle.TaxesConfig).indexOf(key)]}</td>
                            <td>
                                {#if Object.keys(dataSingle.TaxesConfig).indexOf(key) < 2}
                                    {value / MONEY_DIVIDER}
                                {:else}
                                    {value} %
                                {/if}
                            </td>
                        </tr>
                    {/each}

                    </tbody>
                </table>
            {/if}

        {/if}
    </div>
</main>

<style>
    #data-single-buttons-block {
        margin: 20px 0;
    }

    #paid-insurance-save-button {
        margin-left: 0;
    }

    .paid-insurance-input {
        width: 60px;
    }

    label.paid-insurance-label {
        width: auto;
    }
</style>