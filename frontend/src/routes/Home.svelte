<script lang="ts">
  import { GenerateDeclarationOne } from "../../wailsjs/go/main/App.js";
  import { SaveConfig } from "../../wailsjs/go/main/App.js";
  import { LoadConfig } from "../../wailsjs/go/main/App.js";

  let res = ""
  let message = ""

  function decl1(): void {
    GenerateDeclarationOne(message).then((result) => (res = result));
  }

  function save_config(): void {
    SaveConfig().then((result) => (res = result));
  }

  function load_config(): void {
    LoadConfig().then(function (result) {
        res = result.FirstName + " " + result.LastName
    });
  }
</script>

<main>
  <div class="input-box" id="input-box">
    <input autocomplete="off" class="input" id="name" bind:value={message} type="text" />
    <button class="btn" on:click={decl1}>Decl1</button>
     <button class="btn" on:click={save_config}>Save Config</button>
     <button class="btn" on:click={load_config}>Load Config</button>
    <div>{res}</div>

    <a href="#/page">Page</a>
  </div>
</main>

<style>
  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

</style>
