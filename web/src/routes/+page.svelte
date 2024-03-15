<script>
    import { signIn, signOut } from '@auth/sveltekit/client'
    import { page } from '$app/stores'
  </script>
  
  {#if $page.data.session?.user}
    <p>Signed in as {$page.data.session.user.email}</p>
    <p>Session data: {$page.data.session.user}</p>
    <button on:click={signOut}>Sign out</button>
    <img src={$page.data.session.user.image} />
  {:else}
    <p>Not signed in.</p>
    <button on:click={() => signIn('github')}>Sign in</button>
  {/if}