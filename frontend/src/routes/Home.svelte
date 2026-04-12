<script lang="ts">
  import { Save, BookText, CircleCheck, CircleAlert, ArrowBigDown, ArrowBigUp } from 'lucide-svelte';
  import { 
    LoadUserConfig, 
    LoadSettingsConfig, 
    LoadTaxesConfig, 
    GenerateDeclarationOne,
    SaveIncomeForm
   } from "../../wailsjs/go/main/App.js";


  import {BrowserOpenURL} from "../wailsjs/runtime";

  import { onMount } from 'svelte';

	onMount(async () => {
    await load_configs()
	});

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
  let configUser
  let configSettings
  let configTaxes

  let res = ""
  let data // TODO

  async function load_configs(): Promise<void> {
    const [user, settings, taxes] = await Promise.all([
      LoadUserConfig(),
      LoadSettingsConfig(),
      LoadTaxesConfig()
    ]);
    configUser = user;
    configSettings = settings;
    configTaxes = taxes;
  }

  function setMinIncome(): void {
    form.TaxedIncome = String(configTaxes.MinInsuranceIncomeCents / configTaxes.Divider);
  }

  function setMaxIncome(): void {
    form.TaxedIncome = String(configTaxes.MaxInsuranceIncomeCents / configTaxes.Divider);
  }

  function saveIncome(): void {
    if (form.WorkDaysSickLeave > form.WorkDaysTotal) {
      alert("Дните в болничен не могат да надвишават общите работни дни.");
      return;
    }
   
    if (form.DayStart > 0 && form.DayEnd > 0 && form.DayStart >= form.DayEnd) {
      alert("Началният ден трябва да бъде преди крайния ден.");
      return;
    }

    if (form.MonthIncome === "" || form.TaxedIncome === "" || form.Year === 0) {
      alert("Моля, попълнете всички задължителни полета.");
      return;
    }

    if (
      form.TaxedIncome < configTaxes.MinInsuranceIncomeCents / configTaxes.Divider 
      || form.TaxedIncome > configTaxes.MaxInsuranceIncomeCents / configTaxes.Divider
     ) {
      alert(`Осигурителният доход трябва да бъде между ${configTaxes.MinInsuranceIncomeCents / configTaxes.Divider} и ${configTaxes.MaxInsuranceIncomeCents / configTaxes.Divider}.`);
      return;
    }

    // TODO: validation

    let formToSave = {
      Month: parseInt(form.Month),
      Year: form.Year,
      MonthIncomeCents: parseFloat(form.MonthIncome) * configTaxes.Divider, 
      TaxedIncomeCents: parseFloat(form.TaxedIncome) * configTaxes.Divider,
      DayStart: form.DayStart,
      DayEnd: form.DayEnd,
      WorkDaysTotal: form.WorkDaysTotal
    }
   
    SaveIncomeForm(formToSave).then((result) => {
      data = result; // TODO: show success message, or error
    });
  }

  function decl1(): void {
    GenerateDeclarationOne(parseInt(form.Month), form.Year).then((result) => (res = result));
  }
</script>

<main>
  <div class="input-box" id="input-box">
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
      <input class="input" min="2026" required id="Year" type="number" bind:value={form.Year} />
    </div>
  </div>

  <div class="form-row">
    <div class="form-group">
      <label for="MonthIncome">Доход за месец</label>
      <input class="input" min="0" required id="MonthIncome" type="text" bind:value={form.MonthIncome} />
    </div>
    </div>

    <div class="form-row">
    <div class="form-group">
      <label for="TaxedIncome">Осигурителен доход</label>
      <input class="input" min="0" required id="TaxedIncome" type="text" bind:value={form.TaxedIncome} />
    
    <button class="btn btn-small" on:click={setMinIncome}>
      <span><ArrowBigDown color="#444" size="16" /> Мин</span>
    </button>

    <button class="btn btn-small" on:click={setMaxIncome}>
      <span><ArrowBigUp color="#444" size="16" /> Макс</span>
    </button>
    </div>
  </div>

  <div class="form-row">
    <div class="form-group">
      <label for="WorkDaysTotal">Общо работни дни</label>
      <input class="input" min="0" max="31" id="WorkDaysTotal" type="number" bind:value={form.WorkDaysTotal} />

      <!-- TODO: open url-->
      <div class="info">Провери работни дни: <a 
      href="#" 
      on:click={() => BrowserOpenURL('https://kik-info.com/spravochnik/calendar/' + form.Year)}
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
      }} />

      <label class='checkbox-label' for="ShowSickLeave">Бях в болничен</label>
    </div>
  </div>

   <fieldset id="SickLeaveFieldset" style="display: none;">
    <div class="form-row">
      <div class="form-group">
        <label for="WorkDaysSickLeave">Дни в болничен</label>
        <input class="input" min="0" max="31" id="WorkDaysSickLeave" type="number" bind:value={form.WorkDaysSickLeave} />
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
      }} />

      <label class='checkbox-label' for="AddStartOrEnd">Приключвам или започвам дейност този месец</label>
    </div>
  </div>

  <fieldset id="AddStartOrEndFieldset" style="display: none;">
  <div class="form-row">
    <div class="form-group">
      <label for="DayStart">Начален ден на дейност</label>
      <input class="input" id="DayStart" type="number" min="0" max="31" bind:value={form.DayStart} />
      <div class="info">Само ако започваш дейност през този месец</div> <!--TODO: check, Tsveta! -->
    </div>

    </div>  

    <div class="form-row">
    <div class="form-group">
      <label for="DayEnd">Краен ден на дейност</label>
      <input class="input" id="DayEnd" type="number" min="0" max="31" bind:value={form.DayEnd} />
      <div class="info">Само ако приключваш дейност през този месец</div> <!--TODO: check, Tsveta! -->
    </div>
    </div>
    </fieldset> 
  
<div class="form-row">
    <div class="form-group submit-group">
   <button class="btn btn-large" on:click={saveIncome}>
    <span><Save color="#444" size="20" /> Запази</span>
  </button>
</div>
</div>
    {#if data}
    <div class="alert success">
    <CircleCheck color="#748733" size="20" />
      {data}
      <!--
      TODO: taxes and insurances calculated, 
      3 months advance payment tax,
      NET income, 
      -->
    </div>
    
     <div>
      <button class="btn" on:click={decl1}>
        <span><BookText color="#444" size="20" /> Генерирай Декларация 1</span> 
      </button>
    </div>

    {/if}
    <!-- TODO
    <button class="btn">
      <span><BookText color="#444" size="20" /> Генерирай Декларация 6</span> 
    </button>
    <div></div>-->
  </div>

<div id="declarations-schedule">
 <h2>Подаваме в НАП</h2>
  <ul>
    <li>Декларация 1 за дължими осигуровки – всеки месец от 25-то число на следващия месец;</li>
    <li>Декларация 6 за дължими осигурителни вноски – до 30.04 на следващата календарна година;</li>
    <li>Декларация по чл. 55 от ЗДДФЛ – за първите три тримесечия, до края на следващия месец, следващ тримесечието;</li>
    <li>Декларация по чл. 50 от ЗДДФЛ – до 30.04 на следващата календарна година;</li>
    <li>При ДДС регистрация – ежемесечни ДДС декларации до 15-то число на следващия месец.</li>
  </ul>

  </div>
</main>

<style>
  #declarations-schedule {
    padding:0 20px;
    text-align: left;
  }
</style>
