<script lang="ts">
  import {LoadAllIncomeData, DeleteData,GenerateDeclarationOne} from "../../wailsjs/go/main/App.js";
  import {BookText} from 'lucide-svelte';
  import { onMount } from 'svelte';

  let data;
  let res = ""

	onMount(async () => {
    await load_data()
  });

  async function load_data(): Promise<void> {
    LoadAllIncomeData().then(function (result) {
        data = result
    });
	}

  async function delete_data(month: number, year: number): Promise<void> {
    if (!confirm("Сигурен ли си, че искаш да изтриеш тези данни?")) return;
    DeleteData(month, year).then(function (result) {
        console.log(result) // TODO: to alert
        load_data()
    });
  }

  function decl1(month, year): void {
    GenerateDeclarationOne(month, year).then(function(result) {
      console.log(result)
      // TODO: show success message, or error
      res = result
    });
  }
</script>

<main>
<div class="input-box" id="input-box">
  <h2>Въведени данни</h2>

  <table class="table">
    <thead>
      <tr>
        <th>Месец</th>
        <th>Доход</th>
        <th>Осигурителен доход</th>
        <th>Д.1</th>
        <th></th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      {#if data}
        {#each data as row}
          <tr>
            <td>{row.Month} / {row.Year}</td>
            <td>{row.MonthIncomeCents / 100}</td>
            <td>{row.TaxedIncomeCents / 100}</td>
             <td>
               <button class="declaration-button" on:click={() => decl1(row.Month, row.Year)}>
                <BookText color="#444" size="20" />
              </button>
             </td>
             <td><!-- TODO: show all data -->></td>
             <td>
              <button class="delete-button" on:click={() => delete_data(row.Month, row.Year)}>
                X
              </button>
              </td>
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>

</main>

<style>
  .delete-button {
    background-color: transparent;
    border: none;
    color: #923a3a;
    cursor: pointer;
  }

  .declaration-button {
    background-color: transparent;
    border: none;
    color: #444;
    cursor: pointer;
  }
</style>