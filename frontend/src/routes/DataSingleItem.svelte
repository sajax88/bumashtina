<script lang="ts">
    import {GenerateDeclarationOne, LoadIncomeDataForMonth} from "../../wailsjs/go/main/App.js";
    import {onMount} from 'svelte';
    import {BookText} from "lucide-svelte";

    export let params;

    const MONEY_DIVIDER = 100;

    let dataSingle;
    let month = parseInt(params.month);
    let year = parseInt(params.year);

    onMount(async () => {
        await load_data_for_month()
    });

    async function load_data_for_month(): Promise<void> {
        LoadIncomeDataForMonth(month, year).then(function (result) {
            dataSingle = result
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
        <h2>Въведени данни за {month}/{year}</h2>
        {#if dataSingle}
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
                    <td>Осигуровки</td>
                    <td>{dataSingle.SocialSecurityToPayCents / MONEY_DIVIDER}</td>
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

                <!-- TODO: button to see TaxesConfig-->
                </tbody>
            </table>
        {/if}

        <div id="data-single-buttons-block">
            <button class="btn btn-large" id="declaration-one-button" on:click={() => decl1()}>
                <span><BookText color="#444" size="20"/> Генерирай Декларация 1</span>
            </button>
        </div>
    </div>
</main>

<style>
    #data-single-buttons-block {
        margin-top: 20px;
    }
</style>