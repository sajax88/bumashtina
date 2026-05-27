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
    let month = parseInt(params.month);
    let year = parseInt(params.year);
    let taxesLabels;
    let showTaxesConfig = false;
    let paidInsuranceValue = 0;

    onMount(async () => {
        await load_data_for_month()
    });

    async function load_data_for_month(): Promise<void> {
        LoadIncomeDataForMonth(month, year).then(function (result) {
            dataSingle = result
            paidInsuranceValue = dataSingle.SocialSecurityReallyPaidCents / MONEY_DIVIDER
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
        document.getElementById(`paid-insurance-input`).style.display = "inline-block";
        document.getElementById(`paid-insurance-save-button`).style.display = "inline-block";

        document.getElementById(`paid-insurance-value`).style.display = "none";
    }

    function savePaidInsurance() {
        document.getElementById(`paid-insurance-input`).style.display = "none";
        document.getElementById(`paid-insurance-save-button`).style.display = "none";

        document.getElementById(`paid-insurance-value`).style.display = "inline";

        dataSingle.SocialSecurityReallyPaidCents = parseInt(Math.ceil(parseFloat(paidInsuranceValue) * MONEY_DIVIDER));
        UpdateForm(dataSingle).then(function (result) {
           console.log(result) // TODO
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

                <!-- TODO: not populated yet, re-check -->
                <tr>
                    <td>Приблизителен данък (авансов данък се изчислява за 3 месеца)</td>
                    <td colspan="2">{dataSingle.TaxesToPayCents / MONEY_DIVIDER}</td>
                </tr>
                <tr>
                    <td>Изчислени осигуровки</td>
                    <td colspan="2">{dataSingle.SocialSecurityToPayCents / MONEY_DIVIDER}</td>
                </tr>
                <tr>
                    <td>Платени осигуровки</td>
                    <td>
                        <span id="paid-insurance-value">{dataSingle.SocialSecurityReallyPaidCents / MONEY_DIVIDER}</span>

                        <input style="display: none" id="paid-insurance-input" type="text"
                               bind:value={paidInsuranceValue}
                        />
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
                                <Check color="#c6d498" size="20"/>
                            {:else}
                                <Ban color="#edb18a" size="20"/>
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

    #paid-insurance-input {
        width:60px;
    }
</style>