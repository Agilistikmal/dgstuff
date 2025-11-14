<script>
  import Footer from './lib/components/Footer.svelte';
  import Navbar from './lib/components/Navbar.svelte'
</script>

<main>
  <Navbar />
  <div class="min-h-screen">
    <slot />
  </div>
  <Footer />
</main>
