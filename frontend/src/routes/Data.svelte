<script lang="ts">
    import {DeleteData, LoadAllIncomeData, UpdateForm} from "../../wailsjs/go/main/App.js";
    import {ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight, Trash2, View, Edit, Save} from 'lucide-svelte';
    import {onMount} from 'svelte';
    import {main} from "../../wailsjs/go/models";
    import IncomeForm = main.IncomeForm;

    const MONEY_DIVIDER = 100;

    let data;
    let currentPage = 1;
    let itemsPerPage = 12;

    $: totalPages = data ? Math.ceil(data.length / itemsPerPage) : 0;
    $: paginatedData = data ? data.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage) : [];

    onMount(async () => {
        await loadData()
    });

    async function loadData(): Promise<void> {
        LoadAllIncomeData().then(function (result) {
            data = result
            currentPage = 1;
        });
    }

    function changePage(page: number): void {
        if (page >= 1 && page <= totalPages) {
            currentPage = page;
        }
    }

    async function deleteData(month: number, year: number): Promise<void> {
        if (!confirm("Сигурни ли сте, че искате да изтриете данните за " + month + "/" + year + "? Изтритите данни не се възстановяват!")) return;
        DeleteData(month, year).then(function (result) {
            console.log(result) // TODO: to alert, some common error message
            loadData()
        });
    }

    function showPaidTaxesEditInput(month: number, year:number) {
        let id = `${month}${year}`
        document.getElementById('paid-taxes-input-' + id).style.display = "inline-block";
        document.getElementById('paid-taxes-save-button-' + id).style.display = "inline-block";

        document.getElementById('paid-taxes-' + id).style.display = "none";
    }

    function savePaidTaxes(row:IncomeForm) {
        let id = `${row.Month}${row.Year}`
        let amount = document.getElementById('paid-taxes-input-' + id).value

        document.getElementById('paid-taxes-input-' + id).style.display = "none";
        document.getElementById('paid-taxes-save-button-' + id).style.display = "none";

        document.getElementById('paid-taxes-' + id).style.display = "inline";
        document.getElementById('paid-taxes-' + id).innerHTML = amount

        row.TaxesReallyPaidCents = parseInt(Math.ceil(parseFloat(amount) * MONEY_DIVIDER));
        UpdateForm(row).then(function (result) {
            console.log(result) // TODO
        });
    }
</script>

<main>
    <div class="input-box" id="input-box">
        <h2>Въведени данни</h2>

        {#if data && data.length > 0 && totalPages > 1}
            <div class="pagination">
                <button class="pagination-button" on:click={() => changePage(1)} disabled={currentPage === 1}>
                    <ChevronsLeft size="20"/>
                </button>
                <button class="pagination-button" on:click={() => changePage(currentPage - 1)}
                        disabled={currentPage === 1}>
                    <ChevronLeft size="20"/>
                </button>
                <span class="pagination-info">Страница {currentPage} от {totalPages}</span>
                <button class="pagination-button" on:click={() => changePage(currentPage + 1)}
                        disabled={currentPage === totalPages}>
                    <ChevronRight size="20"/>
                </button>
                <button class="pagination-button" on:click={() => changePage(totalPages)}
                        disabled={currentPage === totalPages}>
                    <ChevronsRight size="20"/>
                </button>
            </div>
        {/if}

        <table class="table">
            <thead>
            <tr>
                <th>Месец</th>
                <th>Доход</th>
                <th>Осигурителен доход</th>
                <th>Платени осигуровки</th>
                <th>Платен данък</th>
                <th></th>
                <th></th>
            </tr>
            </thead>
            <tbody>
            {#if paginatedData && paginatedData.length > 0}
                {#each paginatedData as row}

                    {#if row.Month % 3 === 0}
                        <tr class="dark-row">
                            <!-- TODO: roman numbers -->
                            <td>{Math.floor((row.Month + 2) / 3)} тримесечие</td>
                            <td></td>
                            <td></td>
                            <td></td>
                            <td>
                            <span id="paid-taxes-{row.Month}{row.Year}">{row.TaxesReallyPaidCents / MONEY_DIVIDER}</span>
                            <input
                                    style="display: none" id="paid-taxes-input-{row.Month}{row.Year}" type="text"
                                    class="paid-taxes-input" placeholder="0.00"
                                    value="{row.TaxesReallyPaidCents / MONEY_DIVIDER}"
                            />
                            <button style="display: none" class="btn btn-small" id="paid-taxes-save-button-{row.Month}{row.Year}"
                                    on:click="{() => savePaidTaxes(row)}">
                                <Save color="#444" size="20"/>
                            </button>
                            </td>
                            <td class="btn-col">
                                <button class="btn btn-small"  on:click="{() => showPaidTaxesEditInput(row.Month, row.Year)}">
                                    <Edit color="#444" size="20"/>
                                </button>
                            </td>
                            <td></td>
                        </tr>
                    {/if}

                    <tr>
                        <td>{row.Month} / {row.Year}</td>
                        <td>{row.MonthIncomeCents / MONEY_DIVIDER}</td>
                        <td>{row.TaxedIncomeCents / MONEY_DIVIDER}</td>
                        <td>{row.SocialSecurityReallyPaidCents / MONEY_DIVIDER}</td>
                        <td></td>
                        <td class="btn-col">
                            <a href="#/item-single/{row.Year}/{row.Month}">
                                <View color="#444" size="20"/>
                            </a>
                        </td>
                        <td class="btn-col">
                            <button class="delete-button" on:click={() => deleteData(row.Month, row.Year)}>
                                <Trash2 color="#d13708" size="20"/>
                            </button>
                        </td>
                    </tr>
                {/each}
            {/if}
            </tbody>
        </table>

        {#if data && data.length > 0 && totalPages > 1}
            <div class="pagination">
                <button class="pagination-button" on:click={() => changePage(1)} disabled={currentPage === 1}>
                    <ChevronsLeft size="20"/>
                </button>
                <button class="pagination-button" on:click={() => changePage(currentPage - 1)}
                        disabled={currentPage === 1}>
                    <ChevronLeft size="20"/>
                </button>
                <span class="pagination-info">Страница {currentPage} от {totalPages}</span>
                <button class="pagination-button" on:click={() => changePage(currentPage + 1)}
                        disabled={currentPage === totalPages}>
                    <ChevronRight size="20"/>
                </button>
                <button class="pagination-button" on:click={() => changePage(totalPages)}
                        disabled={currentPage === totalPages}>
                    <ChevronsRight size="20"/>
                </button>
            </div>
        {/if}
    </div>

</main>

<style>
    .delete-button {
        background-color: transparent;
        border: none;
        color: #923a3a;
        cursor: pointer;
    }

    .btn-col {
        text-align: center;
    }

    .pagination {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 10px;
        margin: 20px 0;
    }

    .pagination-button {
        padding: 8px 16px;
        background-color: #f0f0f0;
        border: 1px solid #ccc;
        border-radius: 4px;
        cursor: pointer;
        font-size: 14px;
    }

    .pagination-button:hover:not(:disabled) {
        background-color: #e0e0e0;
    }

    .pagination-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .pagination-info {
        font-size: 14px;
        color: #444;
        font-weight: bold;
    }

    .paid-taxes-input {
        width: 100px;
        padding: 6px;
        border: 1px solid #ccc;
        border-radius: 4px;
    }
</style>