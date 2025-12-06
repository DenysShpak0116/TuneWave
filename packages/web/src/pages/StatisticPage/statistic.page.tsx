import { StatisticBlock } from "@modules/StatisticBlock"
import { MainLayout } from "@ui/layout/main-layout"
import { FC, useState } from "react"
import { Nav, NavButton } from "./statistic.page.style";
import { TopListens } from "@modules/TopListens";
import { Combination } from "@modules/Combinations";

export const StatisticPage: FC = () => {
    const [activeTab, setActiveTab] = useState<"listens" | "top" | "comb">("listens");

    return (
        <MainLayout>
            <Nav>
                <NavButton active={activeTab === "listens"} onClick={() => setActiveTab("listens")}>
                    Статистика прослуховування
                </NavButton>
                <NavButton active={activeTab === "top"} onClick={() => setActiveTab("top")}>
                    Топ треків
                </NavButton>
                <NavButton active={activeTab === "comb"} onClick={() => setActiveTab("comb")}>
                    Комбінації
                </NavButton>
            </Nav>

            {activeTab === "listens" && <StatisticBlock />}
            {activeTab === "top" && <TopListens />}
            {activeTab === "comb" && <Combination />}
        </MainLayout>
    );
};