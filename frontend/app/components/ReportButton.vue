<template>
    <div>
        <Dialog v-model:open="open">
            <DialogTrigger as-child>
                <Button variant="outline"
                    class="dark transition-all duration-200 hover:scale-105 hover:bg-primary hover:text-primary-foreground hover:border-primary active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2">
                    <Icon name="lucide:file-text" class="w-4 h-4 mr-2" aria-hidden="true" />
                    Submit Accessibility Report
                </Button>
            </DialogTrigger>
            <DialogContent
                class="bg-card dark text-card-foreground border-border sm:max-w-[625px] max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle class="text-2xl font-bold flex items-center gap-2">
                        <Icon name="lucide:clipboard-pen" class="w-6 h-6 text-primary" aria-hidden="true" />
                        Submit Accessibility Report
                    </DialogTitle>
                    <DialogDescription class="text-muted-foreground">
                        Please provide information about the game's accessibility features to help others.
                    </DialogDescription>
                </DialogHeader>
                <form @submit.prevent="onSubmit" class="space-y-4">
                    <div class="p-1">
                        <FormField v-slot="{ componentField }" name="colorBlind">
                            <FormItem>
                                <FormLabel for="colorBlind" class="text-right">
                                    Color Blind Support
                                </FormLabel>
                                <Select v-bind="componentField" name="colorBlind" class="col-span-3">
                                    <FormControl>
                                        <SelectTrigger class="dark">
                                            <SelectValue placeholder="Select option" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent class="dark">
                                        <SelectItem value="false">There is no support</SelectItem>
                                        <SelectItem value="unknown">I don't know</SelectItem>
                                        <SelectItem value="limited">
                                            There is limited support
                                        </SelectItem>
                                        <SelectItem value="true">There is full support</SelectItem>
                                    </SelectContent>
                                </Select>
                                <ErrorMessage name="colorBlind" class="text-sm text-red-500" />
                            </FormItem>
                        </FormField>
                    </div>

                    <!-- Closed Caption Support Field -->
                    <div class="p-1">
                        <FormField v-slot="{ componentField }" name="closedCaptions">
                            <FormItem>
                                <FormLabel for="closedCaptions" class="text-right">
                                    Closed Caption Support
                                </FormLabel>
                                <Select v-bind="componentField" name="closedCaptions" class="col-span-3">
                                    <FormControl>
                                        <SelectTrigger class="dark">
                                            <SelectValue placeholder="Select option" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent class="dark">
                                        <SelectItem value="false">There is no support</SelectItem>
                                        <SelectItem value="unknown">I don't know</SelectItem>
                                        <SelectItem value="limited">
                                            There is limited support
                                        </SelectItem>
                                        <SelectItem value="true">There is full support</SelectItem>
                                    </SelectContent>
                                </Select>
                                <ErrorMessage name="closedCaptions" class="text-sm text-red-500" />
                            </FormItem>
                        </FormField>
                    </div>

                    <!-- Full Controller Support Field -->
                    <!-- Full Controller Support Field -->
                    <div class="p-1">
                        <FormField v-slot="{ componentField }" name="controllerSupport">
                            <FormItem>
                                <FormLabel for="controllerSupport" class="text-right">
                                    Full Controller Support
                                </FormLabel>
                                <Select v-bind="componentField" name="controllerSupport" class="col-span-3">
                                    <FormControl>
                                        <SelectTrigger class="dark">
                                            <SelectValue placeholder="Select option" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent class="dark">
                                        <SelectItem value="false">There is no support</SelectItem>
                                        <SelectItem value="unknown">I don't know</SelectItem>
                                        <SelectItem value="limited">
                                            There is limited support
                                        </SelectItem>
                                        <SelectItem value="true">There is full support</SelectItem>
                                    </SelectContent>
                                </Select>
                                <ErrorMessage name="controllerSupport" class="text-sm text-red-500" />
                            </FormItem>
                        </FormField>
                    </div>

                    <!-- Controller Remapping Field -->
                    <div class="p-1">
                        <FormField v-slot="{ componentField }" name="controllerRemapping">
                            <FormItem>
                                <FormLabel for="controllerRemapping" class="text-right">
                                    Controller Remapping
                                </FormLabel>
                                <Select v-bind="componentField" name="controllerRemapping" class="col-span-3">
                                    <FormControl>
                                        <SelectTrigger class="dark">
                                            <SelectValue placeholder="Select option" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent class="dark">
                                        <SelectItem value="false">There is no support</SelectItem>
                                        <SelectItem value="unknown">I don't know</SelectItem>
                                        <SelectItem value="limited">
                                            There is limited support
                                        </SelectItem>
                                        <SelectItem value="true">There is full support</SelectItem>
                                    </SelectContent>
                                </Select>
                                <ErrorMessage name="controllerRemapping" class="text-sm text-red-500" />
                            </FormItem>
                        </FormField>
                    </div>

                    <!-- Score Field -->
                    <div class="p-1">
                        <FormField v-slot="{ componentField }" name="score">
                            <FormItem>
                                <FormLabel for="score" class="text-right">
                                    Score (out of 10)
                                </FormLabel>
                                <FormControl>
                                    <Input v-bind="componentField" type="number" min="0" max="10"
                                        class="col-span-3 dark" />
                                </FormControl>
                                <ErrorMessage name="score" class="text-sm text-red-500" />
                            </FormItem>
                        </FormField>
                    </div>

                    <!-- Additional Comments Field -->
                    <div class="p-1">
                        <FormField v-slot="{ componentField }" name="report">
                            <FormItem>
                                <FormLabel for="report" class="text-right">
                                    Additional Comments
                                </FormLabel>
                                <FormControl>
                                    <Textarea v-bind="componentField" placeholder="Optional report for others"
                                        class="col-span-3 dark" />
                                </FormControl>
                                <ErrorMessage name="report" class="text-red-500" />
                            </FormItem>
                        </FormField>
                    </div>

                    <DialogFooter class="flex-col sm:flex-row gap-3">
                        <div v-if="uniqueReport" role="alert" aria-live="assertive"
                            class="flex items-center gap-2 text-destructive text-sm animate-shake">
                            <Icon name="lucide:alert-circle" class="w-4 h-4" aria-hidden="true" />
                            <span>You have already submitted a report for this game</span>
                        </div>
                        <Button type="submit" :disabled="isSubmitting" :aria-busy="isSubmitting"
                            class="transition-all dark duration-200 hover:scale-105 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100">
                            <span v-if="!isSubmitting" class="flex items-center gap-2">
                                <Icon name="lucide:send" class="w-4 h-4" aria-hidden="true" />
                                Submit Report
                            </span>
                            <span v-else class="flex items-center gap-2">
                                <Icon name="lucide:loader-2" class="w-4 h-4 animate-spin" aria-hidden="true" />
                                Submitting...
                            </span>
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    </div>
</template>

<script lang="ts" setup>
import { ErrorMessage, useForm } from "vee-validate";
import * as yup from "yup";

// State management
const uniqueReport = ref(false);
const open = ref(false);
const isSubmitting = ref(false);

const props = defineProps<{
    game: number;
}>();
const emit = defineEmits(["submit"]);

// Define Yup validation schema
const schema = yup.object({
    colorBlind: yup.string().required("Please select an option"),
    closedCaptions: yup.string().required("Please select an option"),
    controllerSupport: yup.string().required("Please select an option"),
    controllerRemapping: yup.string().required("Please select an option"),
    score: yup
        .number()
        .typeError("Score must be a number")
        .required("Please enter a score")
        .min(0, "Score must be at least 0")
        .max(10, "Score cannot exceed 10"),
    report: yup.string().optional(),
});

// Initialize useForm with the validation schema
const form = useForm({
    validationSchema: schema,
});


// Handle form submission
const onSubmit = form.handleSubmit(async (values) => {
    isSubmitting.value = true;
    uniqueReport.value = false;

    console.log("Form Values:", values);
    const reportSubmit = {
        game_id: props.game,
        closed_captions: values.closedCaptions,
        color_blind: values.colorBlind,
        full_controller_support: values.controllerSupport,
        controller_remapping: values.controllerRemapping,
        score: values.score,
        report: values.report,
    };

    await $fetch("/api/report/report", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(reportSubmit),
    })
        .catch((error) => {
            console.log(error.statusMessage);
            if (error.statusCode === 409) {
                uniqueReport.value = true;
            }
        })
        .then(() => {
            if (!uniqueReport.value) {
                emit("submit");
                open.value = false;
            }
        })
        .finally(() => {
            isSubmitting.value = false;
        });
});
</script>

<style scoped>
/* Add any component-specific styles here */

/* Example: Highlight invalid fields */
</style>
