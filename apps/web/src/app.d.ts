declare global {
	namespace App {
		interface Locals {
			user: import('$lib/types').User | null;
		}
		interface Error {
			message: string;
			id?: string;
		}
	}
}

export {};
