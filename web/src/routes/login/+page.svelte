<script>
    /** @type {import('./$types').PageData} */
    export let data;
    import Form from "$lib/components/auth_form/auth_form.svelte";
</script>

<div class="container flex h-screen">
    <div class="m-auto w-[400px]">
        <h1 class="font-semibold text-2xl text-center">Login</h1>
        <p class="text-muted-foreground text-center">Unlock your account.</p>
        <div class="h-4" />
        <Form />
        <div class="h-4" />
        <div>
            Don't have an account?
            <a href="/signup" class="underline hover:text-primary">
                Signup here</a>.
            </div>
    </div>
</div>


