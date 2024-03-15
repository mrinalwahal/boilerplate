<script>
    /** @type {import('./$types').LayoutData} */
    export let data;
</script>

    <div class="grid grid-cols-2 grid-rows-1 gap-0">
        <div>
            <a href="/">
                <img alt="logo" class="w-20 h-20" />
            </a>
        </div>
        <div class="my-auto">
            <slot />
        </div>
    </div>
