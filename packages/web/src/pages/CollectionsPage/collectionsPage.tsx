import { UserCollections } from "@components/UserCollections/user-collections.component";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { useUserCollections } from "@modules/SelectCollectionModal/hooks/useUserCollections";
import { MainLayout } from "@ui/layout/main-layout";
import { Loader } from "@ui/Loader/loader.component";
import { FC } from "react";

export const CollectionsPage: FC = () => {
    const userId = useAuthStore(state => state.user?.id)
    const { data: collections = [], isLoading: loadCollections } = useUserCollections(userId!);

    if (loadCollections) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        )
    }

    return (
        <MainLayout>
            {collections.length > 0 && (
                <UserCollections
                    collections={collections}
                    title="Ваші колекції" />
            )}
        </MainLayout>
    )
}