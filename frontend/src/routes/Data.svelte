<script lang="ts">
  import {LoadAllIncomeData, DeleteData} from "../../wailsjs/go/main/App.js";

   import { onMount } from 'svelte';

  let data;

	onMount(async () => {
    await load_data()
  });

  async function load_data(): Promise<void> {
    LoadAllIncomeData().then(function (result) {
        data = result
    });
	}

  async function delete_data(month: number, year: number): Promise<void> {
    DeleteData(month, year).then(function (result) {
        console.log(result) // TODO: to alert
        load_data()
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
        <th></th>
        <!-- TODO -->
      </tr>
    </thead>
    <tbody>
      {#if data}
        {#each data as row}
          <tr>
            <td>{row.Month} / {row.Year}</td>
            <td>{row.MonthIncomeCents / 100}</td>
            <td>{row.TaxedIncomeCents / 100}</td>
            <!-- TODO -->
             <td>
              <button class="delete-button" on:click={() => delete_data(row.Month, row.Year)}>
                X
              </button>
              </td>
          </tr>
        {/each}
      {/if}
      <!-- TODO -->
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
</style>