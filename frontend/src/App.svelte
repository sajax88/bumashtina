<script lang="ts">
  import Router from "svelte-spa-router";
  import Home from './routes/Home.svelte';

  import Page from './routes/Page.svelte'
</script>

<Router
    routes={{
        "/": Home,
        "/page": Page,
    }}
/>