<script>
  import { onMount } from "svelte";
  import { StuffApi } from "../lib/api/stuff";
  import Navbar from "../lib/components/Navbar.svelte";
  import Footer from "../lib/components/Footer.svelte";
  import emblaCarouselSvelte from "embla-carousel-svelte";
  import Autoplay from "embla-carousel-autoplay";
  import Icon from "@iconify/svelte";
  import { Link } from "svelte-routing";

  let { slug } = $props();

  let loading = $state(true);
  let error = $state(undefined);

  let stuff = $state(null);
  onMount(async () => {
    try {
      loading = true;
      error = undefined;
      stuff = await StuffApi.getBySlug(slug);
    } catch (error) {
      console.error(error);
      error = error.message;
    } finally {
      loading = false;
    }
  });

  let quantity = $state(1);
</script>

<Navbar />
<div class="min-h-screen container mx-auto p-4 pb-24">
  {#if loading}
    <div class="flex justify-center items-center h-screen">
      <div
        class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-brand"
      ></div>
    </div>
  {/if}

  {#if error}
    <div class="flex justify-center items-center h-screen">
      <span class="text-red-500">{error}</span>
    </div>
  {/if}

  {#if stuff}
    {@const stock = stuff.stock.count}
    <Link to="/" class="flex items-center gap-2 w-max">
      <div class="bg-brand text-white rounded-full p-1">
        <Icon icon="basil:caret-left-outline" width="24" height="24" />
      </div>
      <span class="text-sm">Back to Home</span>
    </Link>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
      <!-- Stuff Details -->
      <div class="col-span-2 p-4 border border-gray-300 rounded-lg">
        <div
          use:emblaCarouselSvelte={{
            options: {
              loop: true,
            },
            plugins: [Autoplay({ delay: 5000 })],
          }}
          class="block overflow-hidden h-max aspect-video rounded-lg bg-gray-300"
        >
          <div class="flex">
            {#each stuff.medias as media}
              <div class="flex-[0_0_100%] min-w-0">
                {#if media.type === "image"}
                  <img
                    src={media.url}
                    alt={media.name}
                    class="w-full h-full object-cover"
                  />
                {:else if media.type === "video"}
                  <video
                    src={media.url}
                    autoplay
                    muted
                    loop
                    class="w-full h-full object-cover"
                  ></video>
                {:else}
                  <div
                    class="w-full h-full bg-gray-300 flex items-center justify-center"
                  >
                    <span class="text-gray-500">No media</span>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        </div>

        <div class="mt-4">
          <h1 class="text-2xl font-bold">{stuff.name}</h1>
          <p class="text-gray-500">{stuff.description}</p>
        </div>
      </div>

      <!-- Form Checkout -->
      <div class="col-span-1 p-4 border border-gray-300 rounded-lg h-max">
        <div class="w-full bg-brand text-white rounded-lg p-2">
          <h2 class="text-lg font-bold text-center">Checkout</h2>
        </div>
        <div class="flex items-center gap-4 justify-between mt-4">
          <div>
            {#if stock > 0}
              <span class="text-xl text-green-500">{stock} items left</span>
            {:else}
              <span class="text-xl text-red-500 font-bold">Out of Stock</span>
            {/if}
          </div>
          <div>
            <span class="text-xl text-gray font-bold">
              {Intl.NumberFormat("id-ID", {
                style: "currency",
                currency: stuff.currency,
              }).format(stuff.price)}
            </span>
          </div>
        </div>
        <form>
          <div class="mt-4">
            <label for="quantity" class="block font-medium text-gray-700"
              >Quantity</label
            >
            <div class="flex items-center gap-2 mt-2">
              <button
                type="button"
                class="bg-brand disabled:bg-gray-300 text-white rounded-lg px-5 py-2"
                onclick={() => (quantity -= 1)}
                disabled={quantity <= 1}
                aria-disabled={quantity <= 1}
              >
                <Icon icon="material-symbols:remove" width="24" height="24" />
              </button>
              <input
                type="number"
                id="quantity"
                name="quantity"
                placeholder="1"
                required
                value={quantity}
                min="1"
                max={stock}
                class="flex-1 border border-gray-300 rounded-lg p-2"
              />
              <button
                type="button"
                class="bg-brand disabled:bg-gray-300 text-white rounded-lg px-5 py-2"
                onclick={() => (quantity += 1)}
                disabled={quantity >= stock}
                aria-disabled={quantity >= stock}
              >
                <Icon icon="material-symbols:add" width="24" height="24" />
              </button>
            </div>
          </div>
          <div class="mt-4">
            <button
              class="border border-brand text-brand rounded-lg px-5 py-2 w-full"
            >
              Buy Now
            </button>
          </div>
        </form>
      </div>
    </div>
  {/if}
</div>
<Footer />
