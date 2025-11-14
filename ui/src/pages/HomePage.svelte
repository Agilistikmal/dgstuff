<script>
  import Navbar from "../lib/components/Navbar.svelte";
  import Footer from "../lib/components/Footer.svelte";
  import { onMount } from "svelte";

  let appInfo = $state(undefined);
  onMount(async () => {
    const response = await fetch("/api/appinfo");
    const data = await response.json();
    appInfo = data;
  });
</script>

<Navbar />
<div class="min-h-screen container mx-auto p-4">
  {#if appInfo == undefined}
    <p>Loading...</p>
  {:else if appInfo == null}
    <p>Error loading app info</p>
  {:else}
    <div class="flex items-end gap-2">
      <h1 class="text-7xl font-black tracking-tighter">{appInfo.name}</h1>
      <p class="text-sm text-brand px-3 py-1 bg-brand/10 rounded-full font-medium">{appInfo.version}</p>
    </div>
    <h2 class="text-3xl font-medium">{appInfo.description}</h2>
    <img src={appInfo.logo_url} alt={appInfo.name} class="w-1/2 rounded-lg" />
  {/if}
</div>
<Footer />
