import { UpdateCollectionForm } from "@modules/UpdateCollectionForm";
import { MainLayout } from "@ui/layout/main-layout";
import { Loader } from "@ui/Loader/loader.component";
import { useGetCollection } from "pages/CollectionPage/hooks/useGetCollection";
import { FC } from "react";
import { useParams } from "react-router-dom";

export const UpdateCollectionPage: FC = () => {
    const { id } = useParams();
    const { data: collection, isLoading } = useGetCollection(id!);

    if (isLoading || !collection) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    return (
        <MainLayout>
            <UpdateCollectionForm collection={collection} />
        </MainLayout>
    )
}