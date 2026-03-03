<script lang="ts">
  import { 
    LoadSettingsConfig, 
    SaveSettingsConfig,
    LoadTaxesConfig,
    LoadTaxesConfigLabels
  } from "../../wailsjs/go/main/App.js";

  import { onMount } from 'svelte';

  import { Save, CircleCheck } from 'lucide-svelte';

	onMount(() => {
    load_config()
	});

  let res = ""
  let config;
  let taxes_config;
  let taxes_labels;

  function saveConfig(): void {
    SaveSettingsConfig(config).then((result) => (res = result));
  }

  function load_config(): void {
    LoadSettingsConfig().then(function (result) {
        config = result
    });
    LoadTaxesConfig().then(function (result) {
        taxes_config = result
    });
    LoadTaxesConfigLabels().then(function (result) {
        taxes_labels = result
    });
  }
</script>

<main>
<div class="input-box" id="input-box">
  <h2>Настройки</h2>

  {#if config}
   <div class="form-row">
    <div class="form-group">
      <label for="IsPregnancyInsuranceEnabled">Общо заболяване и майчинство</label>
      <input class="checkbox" id="IsPregnancyInsuranceEnabled" bind:checked={config.IsPregnancyInsuranceEnabled} type="checkbox" />
    </div>
  </div>

   <div class="form-row">
    <div class="form-group submit-group">
    <button class="btn btn-large" on:click={saveConfig}>
    <span><Save color="#444" size="20" /> Запази</span>
  </button>
</div>
</div>

{#if res}
    <div class="alert success">
    <CircleCheck color="#748733" size="20" />  {res}
    </div>
{/if}

{/if}
   <h2>Данъци и осигуровки</h2>
   {#if taxes_config && taxes_labels}
     <table class="table">
       <tbody>
         {#each Object.entries(taxes_config) as [key, value]}
            {#if key !== "Divider"}
           <tr>
             <td>{taxes_labels[Object.keys(taxes_config).indexOf(key)]}</td>
             <td>
                {#if Object.keys(taxes_config).indexOf(key) < 2}
                  {value} <!-- TODO: divide by Divider -->
                {:else}
                  {value} %
                {/if}
            </td>
           </tr>
           {/if}
         {/each}
       </tbody>
     </table>
   {/if}
   </div> 
</main>

<style>
  .table td:first-child {
    width: 350px;
  }
   .table td:last-child {
    font-weight: bold;
  }
</style>