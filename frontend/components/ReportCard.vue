<template>
    <Card class="dark w-full flex">
        <div class="flex-shrink-0 w-24 bg-zinc-600 flex flex-col items-center justify-center p-2 rounded-l-lg">
            <span class="text-sm font-semibold text-center">Accessibility Score</span>
            <span class="text-2xl font-bold">{{ report.score }}</span>
        </div>
        <div class="flex-grow">
            <CardHeader>
                <CardTitle>{{ report.username }}</CardTitle>
                <CardDescription>
                    {{
                        new Date(report.created_at).toLocaleDateString("en-US", {
                            weekday: "long",

                            month: "long",
                            day: "numeric",
                    })
                    }}
                </CardDescription>
            </CardHeader>
            <CardContent>
                <p>{{ report.report }}</p>
            </CardContent>
        </div>
    </Card>
</template>

<script lang="ts" setup>
const props = defineProps<{
    report: Report;
}>();

const config = useRuntimeConfig();
const baseURL = import.meta.server ? config.apiUrl : config.public.apiUrl;

// Only fetch user if user_id exists
const { data: user } = await useFetch<User>(
    () => props.report?.user_id ? `/user/${props.report.user_id}` : null,
    {
        baseURL,
        key: `report-${props.report?.id}-user-${props.report?.user_id}`,
        immediate: !!props.report?.user_id
    }
);

const report = computed(() => ({
    ...props.report,
    username: user.value?.username || props.report.username
}));
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
