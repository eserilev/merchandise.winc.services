Run the merchandise tool with a single parameter:

  ./merchandise -file ./theFile.csv

The file at theFile.csv should contain campaign entries similar to the following. (first line is column names, other lines are data):

Campaign Name,Campaign Tag/Coupon Code,Brand,Category,Channel,Platform,Coupon Type,Max Uses,Minimum Subtotal,Discount Amount,Max Discount,Start Date,Expiration Date,Violator Copy,CAID
Influencer: Influencer Response: Champagne Chaos,champagnechaos,Winc,Non-Gifting,Influencers,NA,Dynamic Discount,10000,29.95,29.95,30,03/01/2022,12/31/2022,New member offer! Get 4 bottles for $29.95 + shipping's on us.,6654178

This is read in and processed by CreateCampaignJSON into files in a pending folder based on the start date of the campaign. Any entries in the processed file will overwrite data in the pending folder. Once everything is generated, these files are used to generate the contents of the campaign-content folder, overwriting whatever is there. This folder contains a file for every day and each member type, if appropriate. Current types are 0 ("leads" or non-members) and 2 (members), allowing us to target campaigns at leads or members.

For example, a campaign with start date 06/12/2022 will end up writing data to campaign-content/2022/06/12/0/index.json (and potentially campaign-content/2022/06/12/2/index.json)


The campaign files can be uploaded to S3 by in the winc-origin-content-develop bucket (i assume) by UploadCampaignFilesToS3. These files are not currently used. Jon commented out this code in June 2022.

The processed file is then moved to an archive directory.

Once the JSON is created, the .go files in content.winc.services must be updated with the new data, and the content.winc.services application must be updated. This is what actually serves the data. Future enhancements here could include reading the data from files uploaded to s3, or having some other way for content.winc.services to ingest the .json files instead of having to be rebuilt.