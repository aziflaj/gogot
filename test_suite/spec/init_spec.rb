RSpec.describe "init" do
  it "requires a path" do
    output = `gogot init`

    expect(output).to include("Usage: gogot init <path>")
  end

  it "creates the directory if missing" do
    output = `rm -rf ./kewl-projekt; gogot init kewl-projekt`
    
    expect(output).to include("Gogot repo initialized in kewl-projekt/.gogot")    
  end
end