// Ce type reflète la structure exacte renvoyée par votre backend Go
export interface UserApiResponse {
    tenantID: string;
    id: string;
    email: string;
    nom: string;
    prenom: string;
    // J'ajoute 'role' en optionnel au cas où votre API évolue,
    // mais l'exemple JSON fourni n'en avait pas.
    role?: string;
}