<script lang="ts">
    import {DeleteData, GenerateDeclarationOne, LoadAllIncomeData} from "../../wailsjs/go/main/App.js";
    import {BookText, ChevronLeft, ChevronsLeft, ChevronRight, ChevronsRight} from 'lucide-svelte';
    import {onMount} from 'svelte';

    let res;
    let data;
    let currentPage = 1;
    let itemsPerPage = 10;

    $: totalPages = data ? Math.ceil(data.length / itemsPerPage) : 0;
    $: paginatedData = data ? data.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage) : [];

    onMount(async () => {
        await load_data()
    });

    async function load_data(): Promise<void> {
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

    async function delete_data(month: number, year: number): Promise<void> {
        if (!confirm("Сигурен ли си, че искаш да изтриеш тези данни?")) return;
        DeleteData(month, year).then(function (result) {
            console.log(result) // TODO: to alert
            load_data()
        });
    }

    function decl1(month, year): void {
        GenerateDeclarationOne(month, year).then(function (result) {
            console.log(result)
            // TODO: show success message, or error
            res = result
        });
    }
</script>

<main>
    <div class="input-box" id="input-box">
        <h2>Въведени данни</h2>

        {#if data && data.length > 0}
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
                <th></th>
                <th>Д.1</th>
                <th></th>
            </tr>
            </thead>
            <tbody>
            {#if paginatedData && paginatedData.length > 0}
                {#each paginatedData as row}
                <tr>
                        <td>{row.Month} / {row.Year}</td>
                        <td>{row.MonthIncomeCents / 100}</td>
                        <td>{row.TaxedIncomeCents / 100}</td>
                        <td><!-- TODO: show all data --></td>
                        <td>
                            <button class="declaration-button" on:click={() => decl1(row.Month, row.Year)}>
                                <BookText color="#444" size="20"/>
                            </button>
                        </td>
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

        {#if data && data.length > 0}
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

    .declaration-button {
        background-color: transparent;
        border: none;
        color: #444;
        cursor: pointer;
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
</style>