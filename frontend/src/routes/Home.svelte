<script lang="ts">
  import { GenerateDeclarationOne } from "../../wailsjs/go/main/App.js";
  import { LoadUserConfig, LoadSettingsConfig } from "../../wailsjs/go/main/App.js";

  import { onMount } from 'svelte';

	onMount(() => {
    load_configs()
	});

  let res = ""
  let form = {
    Month: String(new Date().getMonth()), // We want a previous month
    Year: new Date().getFullYear(),
    MonthIncome: "", 
    TaxedIncome: "",
    DayStart: 0,
    DayEnd: 0,
    WorkDaysTotal: 0,
  }
  let configUser
  let configSettings

  function load_configs(): void {
    LoadUserConfig().then(function (result) {
        configUser = result
    });
    LoadSettingsConfig().then(function (result) {
        configSettings = result
    });
  }

  function decl1(): void {
    GenerateDeclarationOne().then((result) => (res = result));
  }

  function calculate(): void {
    // TODO: calculate in BE
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


    <button class="btn">Set Min</button>
    <button class="btn">Set Max</button>
    </div>
  </div>

  <div class="form-row">
    <div class="form-group">
      <label for="DayStart">Day Start</label>
      <input class="input" id="DayStart" type="number" bind:value={form.DayStart} />
    </div>
    <div class="form-group">
      <label for="DayEnd">Day End</label>
      <input class="input" id="DayEnd" type="number" bind:value={form.DayEnd} />
    </div>
    <div class="form-group">
      <label for="WorkDaysTotal">Work Days Total</label>
      <!-- TODO: open url-->
      <a href="https://kik-info.com/spravochnik/calendar/{form.Year}" target="_blank">Check</a>
      <input class="input" id="WorkDaysTotal" type="number" bind:value={form.WorkDaysTotal} />
    </div>
  </div>

   <button class="btn" on:click={calculate}>Enter</button>
   
    <div>
      TODO: taxes and insurances calculated, 
      3 months advance payment tax,
      NET income, 
    </div>

    <button class="btn" on:click={decl1}>Generate Declaration 1</button>
    <div>{res}</div>

    <!-- TODO
    <button class="btn">Generate Declaration 6 (yearly)</button>
    <div></div>-->
  </div>
</main>

<style>
  
</style>
