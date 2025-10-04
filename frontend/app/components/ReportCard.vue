<template>
    <Card class="dark max-w-lg overflow-hidden transition-all duration-300 hover:-translate-y-1 hover:shadow-xl group"
        role="article" :aria-label="`Report by ${report.username} with score ${report.score}`">
        <!-- Header with Score -->
        <CardHeader class="">
            <div class="flex items-start justify-between gap-4">
                <div class="flex-1 min-w-0">
                    <CardTitle class="text-lg flex items-center gap-2 mb-2">
                        <Icon name="lucide:user" class="w-4 h-4 text-primary flex-shrink-0" aria-hidden="true" />
                        <span class="truncate">{{ report.username }}</span>
                    </CardTitle>
                    <CardDescription class="flex items-center gap-2">
                        <Icon name="lucide:calendar" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                        <time :datetime="report.created_at" class="text-sm">
                            {{
                                new Date(report.created_at).toLocaleDateString("en-US", {
                                    month: "short",
                                    day: "numeric",
                                    year: "numeric",
                                })
                            }}
                        </time>
                    </CardDescription>
                </div>

                <!-- Score Badge -->
                <div class="flex-shrink- flex dark flex-col items-center justify-center bg-neutral-700 rounded-lg p-2 min-w-[3rem] transition-all duration-300 group-hover:scale-105"
                    role="status" :aria-label="`Accessibility score: ${report.score} out of 10`">
                    <span class="text-xl font-bold text-primary leading-none">{{ report.score }}</span>
                </div>
            </div>
        </CardHeader>

        <!-- Content -->
        <CardContent v-if="report.report" class="-mt-4">
            <div class="relative pl-4 border-l-2 border-primary/20">
                <p class="text-muted-foreground leading-relaxed text-lg">{{ report.report }}</p>
            </div>
        </CardContent>
        <CardContent v-else class="pt-0">
            <p class="text-muted-foreground/60 italic text-sm">No additional comments provided</p>
        </CardContent>
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
