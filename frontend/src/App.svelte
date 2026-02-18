<script lang="ts">
  import Router from "svelte-spa-router";
  import Home from './routes/Home.svelte';

  import Settings from './routes/Settings.svelte'
  import Personal from './routes/Personal.svelte'
</script>

<Router
    routes={{
        "/": Home,
        "/settings": Settings,
        "/personal": Personal,
    }}
/>