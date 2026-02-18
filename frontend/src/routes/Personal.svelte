 <script lang="ts">
  import { SaveUserConfig, LoadUserConfig } from "../../wailsjs/go/main/App.js";
  import { onMount } from 'svelte';

  onMount(() => {
        load_config()
	});

  let res = ""
  let config;

  function save_config(): void {
    SaveUserConfig(config).then((result) => (res = result));
  }

  function load_config(): void {
    LoadUserConfig().then(function (result) {
        config = result
    });
  }
</script>

<main>
  <div class="input-box" id="input-box">
    
  <h1>Personal</h1>
  <a href="#/">Home</a>
 
 {#if config}
      <div class="form-group">
        <label for="FirstName">First Name</label>
        <input class="input" required id="FirstName" bind:value={config.FirstName} type="text" />
      </div>
      <div class="form-group">
        <label for="MiddleName">Middle Name</label>
        <input class="input" id="MiddleName" bind:value={config.MiddleName} type="text" />
      </div>
      <div class="form-group">
        <label for="LastName">Last Name</label>
        <input class="input" required id="LastName" bind:value={config.LastName} type="text" />
      </div>
      <div class="form-group">
        <label for="YearOfBirth">Year of Birth</label>
        <input class="input" required type="text" id="YearOfBirth" bind:value={config.YearOfBirth} />
      </div>
      <div class="form-group">
        <label for="Egn">EGN</label>
        <input class="input" required type="text" id="Egn" bind:value={config.Egn} />
      </div>
      <div class="form-group">
        <label for="Bulstat">Bulstat</label>
        <input class="input" required type="text" id="Bulstat" bind:value={config.Bulstat} />
      </div>
      <button class="btn" on:click={save_config}>Save User Config</button>
    <div>{res}</div>
    {/if}
  </div>
</main>

<style>
  
</style>
