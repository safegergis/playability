<template>
  <div class="min-h-screen flex items-center justify-center">
    <Card class="max-w-md shadow-md rounded-md p-6">
      <CardHeader>
        <CardTitle class="text-2xl font-semibold text-center"
          >User Registration</CardTitle
        >
      </CardHeader>
      <CardContent>
        <Form v-slot="{ meta }" :validation-schema="schema" @submit="onSubmit">
          <div class="mb-4">
            <label for="username" class="block text-gray-700 mb-2"
              >Username</label
            >
            <p class="text-gray-500 text-sm mb">
              This is your public display name
            </p>
            <Field
              id="username"
              name="username"
              type="text"
              placeholder="Enter your username"
              class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
            <ErrorMessage name="username" class="text-red-500 text-sm mt-1" />
          </div>

          <div class="mb-4">
            <label for="email" class="block text-gray-700 mb-2">Email</label>
            <Field
              id="email"
              name="email"
              type="email"
              placeholder="Enter your email"
              class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
            <ErrorMessage name="email" class="text-red-500 text-sm mt-1" />
          </div>

          <div class="mb-4">
            <label for="password" class="block text-gray-700 mb-2"
              >Password</label
            >
            <Field
              id="password"
              name="password"
              type="password"
              placeholder="Enter your password"
              class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
            <ErrorMessage name="password" class="text-red-500 text-sm mt-1" />
          </div>

          <div class="mb-6">
            <label for="confirmPassword" class="block text-gray-700 mb-2"
              >Confirm Password</label
            >
            <Field
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              placeholder="Confirm your password"
              class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
            <ErrorMessage
              name="confirmPassword"
              class="text-red-500 text-sm mt-1"
            />
          </div>

          <Button type="submit" class="w-full" :disabled="!meta.valid">
            Register
          </Button>
        </Form>
      </CardContent>
    </Card>
    <Toaster />
  </div>
</template>

<script lang="ts" setup>
import { Form, Field, ErrorMessage } from "vee-validate";
import * as yup from "yup";
import { useToast } from "@/components/ui/toast";
const { toast } = useToast();

// Define the validation schema using Yup
const schema = yup.object({
  username: yup
    .string()
    .required("Username is required")
    .min(3, "Username must be at least 3 characters"),
  email: yup
    .string()
    .required("Email is required")
    .email("Must be a valid email"),
  password: yup
    .string()
    .required("Password is required")
    .min(6, "Password must be at least 6 characters"),
  confirmPassword: yup
    .string()
    .required("Please confirm your password")
    .oneOf([yup.ref("password")], "Passwords must match"),
});

// Handle form submission
const onSubmit = (
  values: {
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
  },
  { resetForm }: { resetForm: () => void }
) => {
  console.log("Form Submitted:", values);
  // Reset form values after submission
  resetForm();

  toast({
    title: "Registration successful!",
    description: "You have successfully registered.",
  });
  // Optionally, you can show a success message or redirect the user
  // For example:
  // useRouter().push('/login');
  // or
  // alert('Registration successful!');
};
</script>
