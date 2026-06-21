<script lang="ts">
    import {
        GenerateDeclarationOne,
        LoadIncomeDataForMonth,
        LoadTaxesConfigLabels,
        UpdateForm
    } from "../../wailsjs/go/main/App.js";
    import {onMount} from 'svelte';
    import {Ban, BookText, Check, CircleArrowLeft, Edit, Eye, EyeOff, Save} from "lucide-svelte";

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
        PensionPartOneCents: 0,
        PensionPartTwoCents: 0,
        HealthInsuranceCents: 0,
    }

    onMount(async () => {
        await load_data_for_month()
    });

    async function load_data_for_month(): Promise<void> {
        LoadIncomeDataForMonth(month, year).then(function (result) {
            //TODO: separate funds
            dataSingle = result
            socialSecurityParts = dataSingle.SocialSecurityToPayParts
            socialSecurityPaidParts = dataSingle.SocialSecurityReallyPaidParts
            if (dataSingle.SocialSecurityReallyPaidCents) {
                paidInsuranceValues.PensionPartOneCents = socialSecurityPaidParts.PensionPartOneCents / MONEY_DIVIDER
                paidInsuranceValues.PensionPartTwoCents = socialSecurityPaidParts.PensionPartTwoCents / MONEY_DIVIDER
                paidInsuranceValues.HealthInsuranceCents = socialSecurityPaidParts.HealthInsuranceCents / MONEY_DIVIDER

            } else {
                // Set to calculated sums
                paidInsuranceValues.PensionPartOneCents = socialSecurityParts.PensionPartOneCents / MONEY_DIVIDER
                paidInsuranceValues.PensionPartTwoCents = socialSecurityParts.PensionPartTwoCents / MONEY_DIVIDER
                paidInsuranceValues.HealthInsuranceCents = socialSecurityParts.HealthInsuranceCents / MONEY_DIVIDER
            }
        });
        LoadTaxesConfigLabels().then(function (result) {
            taxesLabels = result
        });
    }

    function decl1(): void {
        GenerateDeclarationOne(month, year).then(function (result) {
            console.log(result)
            // TODO: show success message, or error
        });
    }

    function showPaidInsuranceEditInput() {
        document.getElementById(`paid-insurance-inputs`).style.display = "inline-block";
        document.getElementById(`paid-insurance-save-button`).style.display = "inline-block";

        document.getElementById(`paid-insurance-value`).style.display = "none";
    }

    async function savePaidInsurance() {
        document.getElementById(`paid-insurance-inputs`).style.display = "none";
        document.getElementById(`paid-insurance-save-button`).style.display = "none";

        document.getElementById(`paid-insurance-value`).style.display = "inline";

        // TODO: functions to turn into cents and back
        dataSingle.SocialSecurityReallyPaidParts.PensionPartOneCents = parseInt(Math.ceil(parseFloat(paidInsuranceValues.PensionPartOneCents) * MONEY_DIVIDER));
        dataSingle.SocialSecurityReallyPaidParts.PensionPartTwoCents = parseInt(Math.ceil(parseFloat(paidInsuranceValues.PensionPartTwoCents) * MONEY_DIVIDER));
        dataSingle.SocialSecurityReallyPaidParts.HealthInsuranceCents = parseInt(Math.ceil(parseFloat(paidInsuranceValues.HealthInsuranceCents) * MONEY_DIVIDER));
        await UpdateForm(dataSingle).then(async function (result) {
            console.log(result) // TODO
            await load_data_for_month()
        });
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
                    <td colspan="2">{dataSingle.MonthIncomeCents / MONEY_DIVIDER}</td>
                </tr>
                <tr>
                    <td>Осигурителен доход</td>
                    <td colspan="2">{dataSingle.TaxedIncomeCents / MONEY_DIVIDER}</td>
                </tr>
                <tr>
                    <td>Изчислен данък<br></td>
                    <td>{dataSingle.TaxesToPayCents / MONEY_DIVIDER}</td>
                    <td>
                        <small>
                            <!-- TODO: paid insurance if entered, take all this from BE? -->
                            ({dataSingle.MonthIncomeCents / MONEY_DIVIDER} доход
                            - {dataSingle.SocialSecurityToPayCents / MONEY_DIVIDER} осигуровки -
                            {dataSingle.MonthIncomeCents / MONEY_DIVIDER} * 0.25 разходи) * 10%
                            = {dataSingle.TaxesToPayCents / MONEY_DIVIDER} EUR
                        </small>
                    </td>
                    <!-- TODO: % from settings, see CalculateTaxForMonth -->
                </tr>
                <tr>
                    <td>Изчислени осигуровки</td>
                    <td>{dataSingle.SocialSecurityToPayCents / MONEY_DIVIDER}</td> <!-- TODO: details -->
                    <td>
                        <small>
                            {socialSecurityParts.PensionPartOneCents / MONEY_DIVIDER} ДОО +
                            {socialSecurityParts.PensionPartTwoCents / MONEY_DIVIDER} ДЗПО +
                            {socialSecurityParts.HealthInsuranceCents / MONEY_DIVIDER} НЗОК
                            = {dataSingle.SocialSecurityToPayCents / MONEY_DIVIDER} EUR
                        </small>
                        <!-- TODO: % from settings, see CalculateSocialSecurity -->
                    </td>
                </tr>
                <tr>
                    <td>Платени осигуровки</td>
                    <td>
                        <span id="paid-insurance-value">
                            {#if dataSingle.SocialSecurityReallyPaidCents}
                                {socialSecurityPaidParts.PensionPartOneCents / MONEY_DIVIDER} ДОО +
                                {socialSecurityPaidParts.PensionPartTwoCents / MONEY_DIVIDER} ДЗПО +
                                {socialSecurityPaidParts.HealthInsuranceCents / MONEY_DIVIDER} НЗОК =
                                {dataSingle.SocialSecurityReallyPaidCents / MONEY_DIVIDER} EUR
                           {:else}
                                -
                            {/if}
                        </span>

                        <span id="paid-insurance-inputs" style="display: none">
                        <label class="paid-insurance-label">ДОО: </label><input
                                class="paid-insurance-input" type="text"
                                bind:value={paidInsuranceValues.PensionPartOneCents}
                        /><br/>
                        <label class="paid-insurance-label">ДЗПО: </label><input
                                class="paid-insurance-input" type="text"
                                bind:value={paidInsuranceValues.PensionPartTwoCents}
                        /><br/>
                        <label class="paid-insurance-label">НЗОК: </label><input
                                class="paid-insurance-input" type="text"
                                bind:value={paidInsuranceValues.HealthInsuranceCents}
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
            <button class="btn btn-small" on:click={() => showTaxesConfig = !showTaxesConfig}>
                            <span>{#if showTaxesConfig}<EyeOff color="#444" size="16"/>{:else}<Eye color="#444"
                                                                                                   size="16"/>{/if}
                                Настройки, данъци и осигуровки</span>
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