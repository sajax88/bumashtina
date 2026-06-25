<script lang="ts">
    import {ArrowBigDown, ArrowBigUp, BookText, CircleCheck, Save} from 'lucide-svelte';
    import {GenerateDeclarationOne, LoadTaxesConfig, SaveIncomeForm,} from "../../wailsjs/go/main/App.js";
    import {BrowserOpenURL} from "../../wailsjs/runtime";
    import {onMount} from 'svelte';

    import Info from "../components/Info.svelte"
    import AlertsBanner from "../components/AlertsBanner.svelte";
    import DeclarationSix from "../components/DeclarationSix.svelte";
    import Taxes from "../components/Taxes.svelte";
    import {main} from "../../wailsjs/go/models";
    import IncomeForm = main.IncomeForm;
    import {fade} from 'svelte/transition';

    let form = {
        Month: String(new Date().getMonth()), // We want a previous month
        Year: new Date().getFullYear(),
        MonthIncome: "",
        TaxedIncome: "",
        DayStart: 0,
        DayEnd: 0,
        WorkDaysTotal: 0,
        WorkDaysSickLeave: 0,
    }


    let configTaxes;
    let declarationResult;

    const MONEY_DIVIDER = 100;

    onMount(async () => {
        await load_configs()
    });

    async function load_configs(): Promise<void> {
        const [taxes] = await Promise.all([
            LoadTaxesConfig(),
        ]);
        configTaxes = taxes;
    }

    function setMinIncome(): void {
        form.TaxedIncome = String(configTaxes.MinInsuranceIncomeCents / MONEY_DIVIDER);
    }

    function setMaxIncome(): void {
        form.TaxedIncome = String(configTaxes.MaxInsuranceIncomeCents / MONEY_DIVIDER);
    }

    function saveIncome(): void {
        let formToSave = new IncomeForm({
            Month: parseInt(form.Month),
            Year: form.Year,
            MonthIncomeCents: Math.ceil(parseFloat(form.MonthIncome) * MONEY_DIVIDER),
            TaxedIncomeCents: Math.ceil(parseFloat(form.TaxedIncome) * MONEY_DIVIDER),
            DayStart: form.DayStart,
            DayEnd: form.DayEnd,
            WorkDaysTotal: form.WorkDaysTotal,
            WorkDaysSickLeave: form.WorkDaysSickLeave,
        })

        SaveIncomeForm(formToSave).then((result) => {
            declarationResult = result;
            setTimeout(() => {
                declarationResult = ""
            }, 30000)
        });
    }

    function generateDeclarationOne(): void {
        GenerateDeclarationOne(parseInt(form.Month), form.Year);
    }
</script>

<main>
    <AlertsBanner/>

    <div id="home-page-block-right" class="input-box">
        <DeclarationSix/>
        <Taxes/>
    </div>

    <div class="input-box" id="home-page-input-box">

        <h2>Въведи данни за доходи</h2>

        <div class="form-row">
            <div class="form-group">
                <label for="Month">Месец</label>
                <select class="input" required id="Month" bind:value={form.Month}>
                    <option value="1">Януари</option>
                    <option value="2">Февруари</option>
                    <option value="3">Март</option>
                    <option value="4">Април</option>
                    <option value="5">Май</option>
                    <option value="6">Юни</option>
                    <option value="7">Юли</option>
                    <option value="8">Август</option>
                    <option value="9">Септември</option>
                    <option value="10">Октомври</option>
                    <option value="11">Ноември</option>
                    <option value="12">Декември</option>
                </select>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group">
                <label for="Year">Година</label>
                <input class="input" min="2026" required id="Year" type="number" bind:value={form.Year}/>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group">
                <label for="MonthIncome">Доход за месец</label>
                <input class="input" min="0" required id="MonthIncome" type="text" bind:value={form.MonthIncome}/>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group">
                <label for="TaxedIncome">Осигурителен доход</label>
                <input class="input" min="0" required id="TaxedIncome" type="text" bind:value={form.TaxedIncome}/>

                <button class="btn btn-small" on:click={setMinIncome}>
                    <span><ArrowBigDown color="#444" size="16"/> Мин</span>
                </button>

                <button class="btn btn-small" on:click={setMaxIncome}>
                    <span><ArrowBigUp color="#444" size="16"/> Макс</span>
                </button>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group">
                <label for="WorkDaysTotal">Общо работни дни</label>
                <input class="input" min="0" max="31" id="WorkDaysTotal" type="number" bind:value={form.WorkDaysTotal}/>

                <div class="info">Провери работни дни: <a href="#"
                                                          on:click={function(e){
                                                              BrowserOpenURL('https://kik-info.com/spravochnik/calendar/' + form.Year);
                                                              e.preventDefault();
                                                          }}
                >
                    https://kik-info.com/spravochnik/calendar/{form.Year}
                </a></div>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group checkbox-group">
                <input id="ShowSickLeave" type="checkbox" on:change={(e) => {
                    document.getElementById('SickLeaveFieldset').style.display = e.target.checked ? 'block' : 'none';
                    if (!e.target.checked) {
                      form.WorkDaysSickLeave = 0;
                    }
                }}/>

                <label class='checkbox-label' for="ShowSickLeave">Бях в болничен</label>
            </div>
        </div>

        <fieldset id="SickLeaveFieldset" style="display: none;">
            <div class="form-row">
                <div class="form-group">
                    <label for="WorkDaysSickLeave">Дни в болничен</label>
                    <input class="input" min="0" max="31" id="WorkDaysSickLeave" type="number"
                           bind:value={form.WorkDaysSickLeave}/>
                </div>
            </div>
        </fieldset>

        <div class="form-row">
            <div class="form-group checkbox-group">
                <input id="AddStartOrEnd" type="checkbox" on:change={(e) => {
        document.getElementById('AddStartOrEndFieldset').style.display = e.target.checked ? 'block' : 'none';
        if (!e.target.checked) {
          form.DayStart = 0;
          form.DayEnd = 0;
        }
      }}/>

                <label class='checkbox-label' for="AddStartOrEnd">Приключвам или започвам дейност този месец</label>
            </div>
        </div>

        <fieldset id="AddStartOrEndFieldset" style="display: none;">
            <div class="form-row">
                <div class="form-group">
                    <label for="DayStart">Начален ден на дейност</label>
                    <input class="input" id="DayStart" type="number" min="0" max="31" bind:value={form.DayStart}/>
                    <div class="info">Само ако започваш дейност през този месец</div>
                </div>

            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="DayEnd">Краен ден на дейност</label>
                    <input class="input" id="DayEnd" type="number" min="0" max="31" bind:value={form.DayEnd}/>
                    <div class="info">Само ако приключваш дейност през този месец</div>
                </div>
            </div>
        </fieldset>

        <div class="form-row">
            <div class="form-group submit-group">
                <button class="btn btn-large" on:click={saveIncome}>
                    <span><Save color="#444" size="20"/> Добави</span>
                </button>
            </div>
        </div>
        {#if declarationResult}
            <div class="alert success" in:fade={{duration:300}} out:fade={{duration:800}}>
                <CircleCheck color="#748733" size="20"/>
                {declarationResult}
            </div>

            <div>
                <button class="btn" on:click={generateDeclarationOne}>
                    <span><BookText color="#444" size="20"/> Генерирай Декларация 1 за {form.Month}/{form.Year}</span>
                </button>
            </div>

        {/if}
    </div>

    <div class="clearfix"></div>

    <Info/>
</main>

<style>
    #home-page-block-right {
        float: right;
        padding-left: 20px;
        width: 35%;
    }

    #home-page-input-box label {
        width: 160px;
    }

    #home-page-input-box div.info {
        padding-left: 170px;
    }

    #home-page-input-box .checkbox-group {
        padding-left: 170px;
    }

    #home-page-input-box label.checkbox-label {
        width: auto;
    }
</style>
