<script lang="ts">
    import { mount_component } from "svelte/internal";
  import { 
    LoadUserConfig, 
    LoadSettingsConfig, 
    LoadTaxesConfig, 
    GenerateDeclarationOne,
    SaveIncomeForm
   } from "../../wailsjs/go/main/App.js";

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
    WorkDaysTotal: 0
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
    console.log(formToSave) // TODO
    SaveIncomeForm(formToSave).then((result) => {
      data = result; // TODO: show success message, or error
    });
  }

  function decl1(): void {
    GenerateDeclarationOne().then((result) => (res = result));
  }
</script>

<main>
    <a href="#/data">Entered Data</a>
    <a href="#/personal">Personal</a>
    <a href="#/settings">Settings</a>

  <div class="input-box" id="input-box">
  <h2>Enter Income</h2>

  <div class="form-row">
    <div class="form-group">
      <label for="Month">Month</label>
      <select class="input" required id="Month" bind:value={form.Month}>
        <option value="1">January</option>
        <option value="2">February</option>
        <option value="3">March</option>
        <option value="4">April</option>
        <option value="5">May</option>
        <option value="6">June</option>
        <option value="7">July</option>
        <option value="8">August</option>
        <option value="9">September</option>
        <option value="10">October</option>
        <option value="11">November</option>
        <option value="12">December</option>
      </select>
    </div>
    <div class="form-group">
      <label for="Year">Year</label>
      <input class="input" required id="Year" type="number" bind:value={form.Year} />
    </div>
  </div>

  <div class="form-row">
    <div class="form-group">
      <label for="MonthIncome">Month Income</label>
      <input class="input" required id="MonthIncome" type="text" bind:value={form.MonthIncome} />
    </div>
    <div class="form-group">
      <label for="TaxedIncome">Taxed Income</label>
      <input class="input" required id="TaxedIncome" type="text" bind:value={form.TaxedIncome} />

    <button class="btn" on:click={setMinIncome}>Set Min</button>
    <button class="btn" on:click={setMaxIncome}>Set Max</button>
    </div>
  </div>

  <div class="form-row">
   <!-- TODO: button to show/hide these two fields-->
    <div class="form-group">
      <label for="DayStart">Day Start</label>
      <input class="input" id="DayStart" type="number" bind:value={form.DayStart} />
    </div>
    <div class="form-group">
      <label for="DayEnd">Day End</label>
      <input class="input" id="DayEnd" type="number" bind:value={form.DayEnd} />
    </div>
    </div>

  <div class="form-row">
    <div class="form-group">
      <label for="WorkDaysTotal">Work Days Total</label>
      <input class="input" id="WorkDaysTotal" type="number" bind:value={form.WorkDaysTotal} />

      <!-- TODO: open url-->
      <br><small> Check: <a href="https://kik-info.com/spravochnik/calendar/{form.Year}" target="_blank">
        https://kik-info.com/spravochnik/calendar/{form.Year}
      </a></small>
    </div>
  </div>

   <button class="btn btn-large" on:click={saveIncome}>Enter</button>

    {#if data}
    <div>
      {data}
      TODO: taxes and insurances calculated, 
      3 months advance payment tax,
      NET income, 
    </div>
    
     <div>
      <button class="btn" on:click={decl1}>Generate Declaration 1</button>
      <div>{res}</div>
    </div>

    {/if}
    <!-- TODO
    <button class="btn">Generate Declaration 6 (yearly)</button>
    <div></div>-->
  </div>

  <!-- TODO: second column -->
  <p>Декларация 1 за дължими осигуровки – всеки месец от 25-то число на следващия месец;<br>
	Декларация 6 за дължими осигурителни вноски – до 30.04 на следващата календарна година; <br>
	Декларация по чл. 55 от ЗДДФЛ – до края на месеца следващ тримесечието, за което декларацията се подава (само за първите три тримесечия); <br>
	Декларация по чл. 50 от ЗДДФЛ – до 30.04 на следващата календарна година <br>
	При ДДС регистрация – ежемесечни ДДС декларации до 15-то число на следващия месец.</p>
</main>

<style>
  .input-box label {
    display:inline-block;
    width:120px;
  }
</style>
