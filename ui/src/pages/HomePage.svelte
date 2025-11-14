<script>
  import Navbar from "../lib/components/Navbar.svelte";
  import Footer from "../lib/components/Footer.svelte";
  import { AppInfoApi } from "../lib/api/appinfo.js";
  import { StuffApi } from "../lib/api/stuff.js";
  import { onMount } from "svelte";
  import { Link } from "svelte-routing";

  let loading = $state(true);
  let error = $state(undefined);

  let page = $state(1);
  let limit = $state(10);

  let appInfo = $state(null);
  let stuffs = $state(null);
  onMount(async () => {
    try {
      loading = true;
      error = undefined;
      appInfo = await AppInfoApi.get();
      stuffs = await StuffApi.getAll(page, limit);
    } catch (error) {
      console.error(error);
      error = error.message;
    } finally {
      loading = false;
    }
  });
</script>

<Navbar />
<div class="min-h-screen container mx-auto p-4 pb-24">
  {#if loading}
    <p>Loading...</p>
  {/if}
  {#if error}
    <p>Error loading app info</p>
  {/if}

  {#if appInfo}
    <!-- Hero -->
    <div class="min-h-[25vh] flex flex-col justify-center items-start">
      <div class="flex items-end gap-2">
        <h1 class="text-7xl font-black tracking-tighter">{appInfo.name}</h1>
        <p
          class="text-sm text-brand px-3 py-1 bg-brand/10 rounded-full font-medium"
        >
          {appInfo.version}
        </p>
      </div>
      <h2 class="text-3xl font-medium">{appInfo.description}</h2>
      <img
        src={appInfo.logo_url}
        alt={appInfo.name}
        class="h-12 rounded-lg mt-4"
      />
    </div>

    <!-- Stuffs -->
    <div id="stuffs">
      <!-- Category Filters -->
      <div class="flex items-center gap-4">
        <button class="rounded-lg px-4 py-2 bg-brand/10">
          <span>Category A</span>
        </button>
      </div>

      <!-- Search Input -->
      <div class="flex items-center gap-4 mt-4">
        <input
          type="text"
          placeholder="Search"
          class="rounded-lg px-4 py-2 border border-gray-300 w-full"
        />
      </div>

      <!-- List -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
        {#each stuffs?.data as stuff}
          <Link
            to={`/stuff/${stuff.slug}`}
            class="rounded-lg overflow-hidden shadow-sm hover:shadow-xl transition-shadow duration-300"
          >
            <img
              src={StuffApi.getFirstMedia(stuff.medias, "image")}
              alt={stuff.name}
              class="w-full h-48 object-cover bg-gray-300 flex items-center justify-center"
            />
            <div class="p-4">
              <h3 class="text-lg font-medium">{stuff.name}</h3>
              <div class="flex items-center gap-2">
                {#each stuff.categories as category, index}
                  <span class="text-sm text-gray-500">
                    {category.name}
                  </span>
                  {#if index !== stuff.categories.length - 1}
                    <span class="text-xs text-accent">•</span>
                  {/if}
                {/each}
              </div>
              <p class="text-brand font-bold">
                {Intl.NumberFormat("id-ID", {
                  style: "currency",
                  currency: stuff.currency,
                }).format(stuff.price)}
              </p>
            </div>
          </Link>
        {/each}
      </div>

      <!-- Pagination -->
      {#if stuffs?.total_items == 0}
        <div class="flex justify-center mt-8">
          <span class="text-gray-500">No stuffs found</span>
        </div>
      {/if}
      <div class="flex justify-center mt-8">
        {#if stuffs?.has_previous}
          <button class="rounded-lg px-4 py-2 bg-brand/10">
            <span>Previous</span>
          </button>
        {/if}
        <span class="text-sm text-gray-500"
          >Page {stuffs?.page} of {stuffs?.total_pages}</span
        >
        {#if stuffs?.has_next}
          <button class="rounded-lg px-4 py-2 bg-brand/10">
            <span>Next</span>
          </button>
        {/if}
      </div>
    </div>
  {/if}
</div>
<Footer />
