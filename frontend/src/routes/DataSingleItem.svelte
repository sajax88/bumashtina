<script lang="ts">
    import {GenerateDeclarationOne, LoadIncomeDataForMonth, LoadTaxesConfigLabels} from "../../wailsjs/go/main/App.js";
    import {onMount} from 'svelte';
    import {BookText, Eye, EyeOff, CircleArrowLeft} from "lucide-svelte";

    export let params;

    const MONEY_DIVIDER = 100;

    let dataSingle;
    let month = parseInt(params.month);
    let year = parseInt(params.year);
    let taxesLabels;
    let showTaxesConfig = false;

    onMount(async () => {
        await load_data_for_month()
    });

    async function load_data_for_month(): Promise<void> {
        LoadIncomeDataForMonth(month, year).then(function (result) {
            dataSingle = result
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
</script>

<main>
    <div class="input-box" id="input-box">
        <a href="#/data"><CircleArrowLeft color="#444" size="20"/></a>

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
                    <td>{dataSingle.MonthIncomeCents / MONEY_DIVIDER}</td>
                </tr>
                <tr>
                    <td>Осигурителен доход</td>
                    <td>{dataSingle.TaxedIncomeCents / MONEY_DIVIDER}</td>
                </tr>

                <!-- TODO: not populated yet, re-check -->
                <tr>
                    <td>Приблизителен данък (авансов данък се изчислява за 3 месеца)</td>
                    <td>{dataSingle.TaxesToPayCents / MONEY_DIVIDER}</td>
                </tr>
                <tr>
                    <td>Изчислени осигуровки</td>
                    <td>{dataSingle.SocialSecurityToPayCents / MONEY_DIVIDER}</td>
                </tr>
                <tr>
                    <td>Платени осигуровки</td>
                    <td>{dataSingle.SocialSecurityReallyPaidCents / MONEY_DIVIDER}</td>
                </tr>
                <tr>
                    <td>Работни дни</td>
                    <td>{dataSingle.WorkDaysTotal}</td>
                </tr>
                {#if dataSingle.WorkDaysSickLeave}
                    <tr>
                        <td>В болничен</td>
                        <td>{dataSingle.WorkDaysSickLeave}</td>
                    </tr>
                {/if}

                {#if dataSingle.DayStart}
                    <tr>
                        <td>Започване на дейност от ден</td>
                        <td>{dataSingle.DayStart}</td>
                    </tr>
                {/if}

                {#if dataSingle.DayEnd}
                    <tr>
                        <td>Приключване на дейност от ден</td>
                        <td>{dataSingle.DayEnd}</td>
                    </tr>
                {/if}

                <tr>
                    <td colspan="2">
                        <button class="btn btn-small" on:click={() => showTaxesConfig = !showTaxesConfig}>
                            <span>{#if showTaxesConfig}<EyeOff color="#444" size="16"/>{:else}<Eye color="#444" size="16"/>{/if}
                                Настройки, данъци и осигуровки</span>
                        </button>
                    </td>
                </tr>

                {#if showTaxesConfig}
                <tr id="data-single-taxes-config">
                    <td>Данъци и осигуровки за {month}/{year}</td>
                    <td>
                        <table class="table inner-table">
                            <tbody>
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
                    </td>
                </tr>
                {/if}
                </tbody>
            </table>
        {/if}
    </div>
</main>

<style>
    #data-single-buttons-block {
        margin: 20px 0;
    }

    .inner-table th, .inner-table td {
        border:0;
    }
</style>