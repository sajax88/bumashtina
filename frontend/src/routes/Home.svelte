<script lang="ts">
  import { GenerateDeclarationOne } from "../../wailsjs/go/main/App.js";
  import { LoadUserConfig, LoadSettingsConfig } from "../../wailsjs/go/main/App.js";

  import { onMount } from 'svelte';

	onMount(() => {
    load_configs()
	});

  let res = ""
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


    <a href="#/personal">Personal</a>
    <a href="#/settings">Settings</a>

  <div class="input-box" id="input-box">
  <h2>Enter Income</h2>

  <div class="form-row">
    <div class="form-group">
      <label for="Month">Month</label>
      <input class="input" required id="Month" type="number" />
    </div>
    <div class="form-group">
      <label for="Year">Year</label>
      <input class="input" required id="Year" type="number" />
    </div>
  </div>

  <div class="form-row">
    <div class="form-group">
      <label for="MonthIncome">Month Income</label>
      <input class="input" required id="MonthIncome" type="text" />
    </div>
    <div class="form-group">
      <label for="TaxedIncome">Taxed Income</label>
      <input class="input" required id="TaxedIncome" type="text" />
    </div>
  </div>

  <div class="form-row">
    <div class="form-group">
      <label for="DayStart">Day Start</label>
      <input class="input" id="DayStart" type="number" />
    </div>
    <div class="form-group">
      <label for="DayEnd">Day End</label>
      <input class="input" id="DayEnd" type="number" />
    </div>
    <div class="form-group">
      <label for="WorkDaysTotal">Work Days Total</label>
      <input class="input" id="WorkDaysTotal" type="number" />
    </div>
  </div>

   <button class="btn" on:click={calculate}>Enter</button>
    <div>
      TODO: taxes and insurances calculated, 
      3 months advance payment tax,
      NET income, 
    </div>

    <button class="btn" on:click={decl1}>Generate Declaration 1 (monthly)</button>
    <div>{res}</div>

    <!-- TODO
    <button class="btn">Generate Declaration 6 (yearly)</button>
    <div></div>-->
  </div>
</main>

<style>
  
</style>
