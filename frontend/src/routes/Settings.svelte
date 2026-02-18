<script lang="ts">
  import { LoadSettingsConfig, SaveSettingsConfig } from "../../wailsjs/go/main/App.js";

  import { onMount } from 'svelte';

	onMount(() => {
    load_config()
	});

  let res = ""
  let config;

  function save_config(): void {
    SaveSettingsConfig(config).then((result) => (res = result));
  }

  function load_config(): void {
    LoadSettingsConfig().then(function (result) {
        config = result
    });
  }
</script>

<main>
  <h1>Settings</h1>
  <a href="#/">Home</a>

  {#if config}
    <div class="form-group">
      <label for="IsPregnancyInsuranceEnabled">Is Pregnancy Insurance Enabled</label>
      <input class="checkbox" id="IsPregnancyInsuranceEnabled" bind:checked={config.IsPregnancyInsuranceEnabled} type="checkbox" />
    </div>
  {/if}
   <button class="btn" on:click={save_config}>Save Settings</button>
</main>