using JEngine.Singleton;
using pb;
using System.Collections.Generic;
using static pb.MainBattlePrepareResponse;

namespace HotUpdateScripts.ProtoM
{
    public class HallBattleM : SingletonTemplate<HallBattleM>
    {
        private List<IDQuality> dropItems;
        private int index;

        public void AddEventListener()
        {
            UniGameEvent.AddEventListener<MainBattlePrepareResponse>(Sever_MainBattlePrepareResponse);
        }

        public void RemoveEventListener()
        {
            UniGameEvent.RemoveEventListener<MainBattlePrepareResponse>(Sever_MainBattlePrepareResponse);
        }

        public void SendMainBattlePrepareRequest()
        {
            dropItems = null;
            index = 0;
            int job = RoleM.Instance.GetDefaultOpenBoxEquipJob();
            TcpC.C2S_MainBattlePrepare(job);
        }

        private void Sever_MainBattlePrepareResponse(MainBattlePrepareResponse req)
        {
            dropItems = req.Items;
            //UnityEngine.Debug.LogError("Sever_MainBattlePrepareResponse");
            //for (int i = 0; i < req.Items.Count; i++)
            //{
            //    UnityEngine.Debug.LogError(req.Items[i].ID);
            //}
        }

        public IDQuality GetDropItemByID(int id)
        {
            for (int i = 0; i < dropItems.Count; i++)
            {
                if (dropItems[i].ID == id)
                    return dropItems[i];
            }
            return null;
        }

        public IDQuality[] GetDropItem(int totalCount)
        {
            var count = dropItems.Count / totalCount;
            var temp = new IDQuality[count];

            for (int i = 0; i < count; i++)
            {
                temp[i] = dropItems[index * count + i];
            }

            index++;
            return temp;
        }
    }
}