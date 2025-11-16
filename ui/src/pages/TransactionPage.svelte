<script>
  import Navbar from "../lib/components/Navbar.svelte";
  import Footer from "../lib/components/Footer.svelte";
  import { onMount } from "svelte";
  import { TransactionApi } from "../lib/api/transaction";
  import Loading from "../lib/components/Loading.svelte";

  let { id } = $props();

  let loading = $state(true);
  let error = $state(undefined);

  let transaction = $state(null);
  onMount(async () => {
    try {
      loading = true;
      error = undefined;
      const token =
        new URLSearchParams(window.location.search).get("token") || "";
      console.log("Token:", token);
      transaction = await TransactionApi.get(id, token);
      if (!transaction) {
        throw new Error(`Transaction with id ${id} not found`);
      }
    } catch (err) {
      console.error(error);
      error = err.message;
    } finally {
      loading = false;
    }
  });

  const getStatusColor = (prefix = "") => {
    switch (transaction.status) {
      case "pending":
        return `${prefix}-brand`;
      case "success":
        return `${prefix}-green-500`;
      case "failed":
        return `${prefix}-red-500`;
    }
  };

  const refreshPaymentStatus = async () => {
    try {
      loading = true;
      error = undefined;
      transaction = await TransactionApi.get(id);
    } catch (err) {
      console.error(error);
      error = err.message;
    } finally {
      loading = false;
    }
  };

  let copied = $state(false);
  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    copied = true;
    setTimeout(() => {
      copied = false;
    }, 2000);
  };
</script>

<Navbar />
<div class="min-h-screen w-full max-w-screen-sm mx-auto p-4 pb-24">
  {#if loading}
    <Loading />
  {/if}
  {#if error}
    <div>
      <span class="text-red-500">{error}</span>
    </div>
  {/if}
  {#if transaction}
    <div
      class={`flex flex-col items-center gap-4 p-4 border border-gray-300 border-b-4 ${getStatusColor("border-b")} rounded-lg`}
    >
      <div class={`w-full ${getStatusColor("bg")} text-white rounded-lg p-2`}>
        <h1 class="text-lg font-bold text-center">Transaction</h1>
      </div>
      <div class="w-full">
        <div class="flex items-center gap-2 justify-between">
          <h2 class="text-lg font-bold">Transaction Details</h2>
          <span class={`uppercase font-bold ${getStatusColor("text")}`}
            >{transaction.status}</span
          >
        </div>
        <div class="flex items-center gap-2 justify-between text-gray text-sm">
          <p>ID: {transaction.id}</p>
          <p>
            {transaction.email}
          </p>
        </div>
      </div>

      <!-- Stuffs -->
      <table class="w-full border border-collapse border-gray-300">
        <thead class="bg-gray-100 text-left">
          <tr>
            <th class="p-2">Name</th>
            <th class="p-2">Quantity</th>
            <th class="p-2">Price</th>
          </tr>
        </thead>
        <tbody class="text-left text-sm align-top">
          {#each transaction.stuffs as stuff}
            <tr>
              <td class="p-2">{stuff.stuff_name}</td>
              <td class="p-2">{stuff.quantity}</td>
              <td class="p-2"
                >{Intl.NumberFormat("id-ID", {
                  style: "currency",
                  currency: transaction.currency,
                }).format(stuff.total_price)}</td
              >
            </tr>
          {/each}
          <tr class="font-bold border-t border-gray-300">
            <td class="p-2">Total</td>
            <td class="p-2">
              {transaction.stuffs.reduce(
                (acc, stuff) => acc + stuff.quantity,
                0,
              )}
            </td>
            <td class={`p-2 ${getStatusColor("text")}`}>
              {Intl.NumberFormat("id-ID", {
                style: "currency",
                currency: transaction.currency,
              }).format(transaction.amount)}
            </td>
          </tr></tbody
        >
      </table>

      <!-- Pay Button -->
      <div class="w-full">
        <button
          onclick={() => {
            window.open(transaction.payment.url, "_blank");
          }}
          class={`w-full border ${getStatusColor("border")} ${getStatusColor("text")} font-bold rounded-lg p-2 capitalize block text-center cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed`}
          disabled={transaction.status !== "pending"}
        >
          {#if transaction.status === "pending"}
            Pay with {transaction.payment.provider}
          {:else if transaction.status === "success"}
            Payment successful
          {:else if transaction.status === "failed"}
            Payment failed
          {/if}
        </button>

        {#if transaction.status === "pending"}
          <div class="text-center w-full mt-2 space-y-2">
            <span class="text-sm text-gray-500 block"
              >Pay before {new Date(
                transaction.payment.expires_at,
              ).toLocaleString("id-ID", {
                dateStyle: "medium",
                timeStyle: "short",
              })}</span
            >
            <span class="text-sm text-gray-500 block">
              Looks like you already paid?
              <button class="text-brand" onclick={refreshPaymentStatus}
                >Refresh Payment Status</button
              >
            </span>
          </div>
        {/if}
      </div>

      <!-- Stuff Values -->
      <div class="flex items-center gap-2 w-full">
        <hr class="flex-1 border-gray-300 my-2" />
        <span class="text-gray-500">Your Purchased Stuff</span>
        <hr class="flex-1 border-gray-300 my-2" />
      </div>
      <div class="w-full space-y-2">
        {#each transaction.stuffs as stuff}
          <div>
            <h1 class="font-medium">{stuff.stuff_name}</h1>
            <p class="text-sm text-gray-500">
              Separator <span
                class="font-bold text-brand bg-brand/10 px-2 rounded-lg"
                >{stuff.data.separator}</span
              >
            </p>
            <div
              class="mt-2 bg-gray-100 p-2 rounded-lg border border-gray-300 relative"
            >
              <code class="block">{stuff.data.values}</code>
              <button
                class="absolute top-2 right-2 text-sm text-gray-500 hover:text-gray-700"
                onclick={() => {
                  copyToClipboard(stuff.data.values);
                }}
              >
                {copied ? "Copied" : "Copy"}
              </button>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>
<Footer />
