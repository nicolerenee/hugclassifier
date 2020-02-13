package classifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicolerenee/hugclassifier/pkg/classifier"
)

var (
	testDesc1 string = `<p>Get pumped for Perth's first HashiCorp User Group Meeting!</p> <p>In our first meeting we will be discussing topic of using Infrastructure as code using Terraform and Packer.</p> <p>6pm<br/>Arrive, sign in, pre-network high-fives and fist bumps</p> <p>Speakers and Talks:<br/>6:30pm<br/>Bruce Dominguez: Site Reliability Engineer - VGW<br/>Resource "Topic" "The Death of ClickOps: IaC with Packer and Terraform" {<br/>Summary = There is always that one server that everyone is scared to turn off, that everyone SSH's on to for deployments. This talk will show you how we went from a snowflake server to immutable infrastructure defined as code using Packer and Terraform, and deployed as part of our CI/CD pipeline.<br/>}</p> <p>7:00pm<br/>Break</p> <p>7:10pm<br/>Even Zhang: Senior Engineer - VGW<br/>Kick start your project with Infrastructure as Code<br/>summary: When you hear infrastructure as code, you may think of it as code that _just_ manages infrastructure such as your servers, your firewalls etc. But today it means so much more, you can use IaC to provision all your permissions, your git repos, your CICD pipelines, your monitoring setup and even getting pizzas ready.</p> <p>7:40pm<br/>Drinks, food and more high-fives and fist bumps!<br/>-----</p> <p>We look forward to seeing you!</p> <p>Thank you to our sponsors<br/>VGW &amp; New Relic</p>`
	testDesc2 string = `<p>Hello Detroit &amp; Ann Arbor!</p> <p>We are hosting another HashiCorp Meetup in January! This month we are going on a field trip to Ann Arbor. The meeting will be held at the offices of SkySpecs.</p> <p>Please join us for an opportunity to meet other users in the community, discuss implementation details and challenges, or just hang out to hear technical stories and get information.</p> <p>Also looking for community members that would like to present at future events, so let us know if you're interested.</p> <p>Feedback on how we can be better is always welcome.</p> <p>----------------------------------------------------------------------------------------------------</p> <p>Agenda:</p> <p>-Introduction/Meet &amp; Greet</p> <p>-Talk #1 - HashiCorp Overview</p> <p>-Talk #2 - HashiPi Fresh from the oven (evaluating the Hashi product stack on a Raspberry Pi cluster)</p> <p>About SkySpecs:</p> <p>SkySpecs entered the wind energy industry in 2014 as a fully automated drone inspection company. After quickly earning a leading position in the North American and European markets, having conducted over 22,000 offshore and onshore inspections for the largest wind energy owners, we grew our services and capabilities to offer full-wrap operations and maintenance solutions.</p> <p>Our multi-layered solutions include analytics, wind turbine blade expertise, engineering projects, collaborative software to manage and analyze data from multiple sources, digitization of data, planning and consultation on high-cost repair campaigns.</p> <p>About HashiCorp:</p> <p>At HashiCorp we are focused on providing tools that allow organizations to adopt the cloud and securely automate their infrastructure, with the end goal of accelerating engagement with your customers, partners, and employees. Our solutions increase speed to market, reduce cost, and manage risk by leveraging 'infrastructure as code' to provision, secure, connect, and run any application, anywhere.</p>`
	testDesc3 string = `<p>Hi everyone! Join us for the Des Moines Iowa inaugural HUG on Tuesday, December 10th at Exile Brewing. We have an amazing roster of speakers to host our first HUG. Join us for a night of food, networking, and relevant HashiCorp related topics</p> <p>Check the agenda out below!</p> <p>Agenda:<br/>*** 4:00pm - Arrivals, food &amp; beverages, networking.</p> <p>*** 4:30 pm – Opening remarks from our sponsors and organizers.</p> <p>*** 4:45 pm – Bryan Schleisman with Hashicorp will demo "Full Application Stack Automation" leveraging open source Hashicorp solutions like Packer, Terraform, Nomad, Consul, and Vault.</p> <p>*** 5:15 pm - Keith Richter with AppDynamics will Demo "APM Solution" from AppDymanics. See what the next steps are to monitor and manage your applications once they are built.</p> <p>*** 5:40 - Closing remarks, food, and drinks.</p> <p>We look forward to seeing you.</p> <p>Interested in getting involved with this chapter? We are always looking for venues and speakers. Reach out to us at [masked].</p>`
)

func TestDesc1Classify(t *testing.T) {
	r := classifier.Classify(testDesc1)
	assert.ElementsMatch(t, []string{"Terraform", "Packer"}, r)
}
func TestDesc2Classify(t *testing.T) {
	r := classifier.Classify(testDesc2)
	assert.ElementsMatch(t, []string{}, r)
}

func TestDesc3Classify(t *testing.T) {
	r := classifier.Classify(testDesc3)
	assert.ElementsMatch(t, []string{"Consul", "Nomad", "Packer", "Terraform", "Vault"}, r)
}
